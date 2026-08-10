package lambda

import (
	"context"
	stderrors "errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// CheckCache will return a non-empty string if the URL given is in the cache
func CheckCache(query TitleQuery) (string, error) {
	ctx := context.TODO()

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("eu-west-1"))
	if err != nil {
		log.Errorf("could not connect to AWS %v", err)
		return "", err
	}

	// Create DynamoDB client
	svc := dynamodb.NewFromConfig(cfg)

	input := &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"url": &types.AttributeValueMemberS{Value: query.URL},
		},
		TableName: aws.String("urls"),
	}

	// From: https://docs.aws.amazon.com/sdk-for-go-v2/api/service/dynamodb/#Client.GetItem
	result, err := svc.GetItem(ctx, input)
	// Error when fetching
	if err != nil {
		logDynamoDBError(err)
		return "", errors.New(err.Error())
	}

	// Grab the title from the result and return it
	// TODO:
	// optionally update ttl in DB -> frequent stuff gets cached longer

	if result.Item != nil {
		titleAttr, ok := result.Item["title"].(*types.AttributeValueMemberS)
		if !ok {
			return "", errors.New("Cache miss")
		}
		return titleAttr.Value, nil
	}
	return "", errors.New("Cache miss")
}

// CacheAndReturn inserts a successfully found title to cache
func CacheAndReturn(query TitleQuery, title string, err error) (TitleQuery, error) {
	if err != nil {
		query.Title = ""
		return query, err
	}

	ctx := context.TODO()

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("eu-west-1"))
	if err != nil {
		log.Errorf("could not connect to AWS %v", err)
		return TitleQuery{}, err
	}

	// Create DynamoDB client
	svc := dynamodb.NewFromConfig(cfg)
	// create a map for DD
	query.Title = title
	query.Added = time.Now().Unix()
	query.TTL = time.Now().Unix() + 86400 // 24 hours

	log.Infof("Storing TitleQuery: %v", query)

	av, err := attributevalue.MarshalMap(query)
	if err != nil {
		log.Errorf("error marshaling to dynamodb: %v", err)
		return TitleQuery{}, err
	}

	// construct an input that DD can handle
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("urls"),
	}
	// put item in DD
	_, err = svc.PutItem(ctx, input)
	if err != nil {
		logDynamoDBError(err)
	}

	return query, err
}

// logDynamoDBError logs known DynamoDB error types with their code and message,
// falling back to a generic log for anything else.
func logDynamoDBError(err error) {
	var throughputErr *types.ProvisionedThroughputExceededException
	var notFoundErr *types.ResourceNotFoundException
	var limitErr *types.RequestLimitExceeded
	var internalErr *types.InternalServerError

	switch {
	case stderrors.As(err, &throughputErr):
		log.Error("ProvisionedThroughputExceededException", throughputErr.Error())
	case stderrors.As(err, &notFoundErr):
		log.Error("ResourceNotFoundException", notFoundErr.Error())
	case stderrors.As(err, &limitErr):
		log.Error("RequestLimitExceeded", limitErr.Error())
	case stderrors.As(err, &internalErr):
		log.Error("InternalServerError", internalErr.Error())
	default:
		log.Error(err.Error())
	}
}
