package lambda

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// CheckCache will return a non-empty string if the URL given is in the cache
func CheckCache(query TitleQuery) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1")},
	)

	if err != nil {
		log.Errorf("could not connect to AWS %v", err)
		return "", err
	}

	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"url": {
				S: aws.String(query.URL),
			},
		},
		TableName: aws.String("urls"),
	}

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	// From: https://docs.aws.amazon.com/sdk-for-go/api/service/dynamodb/#DynamoDB.GetItem
	result, err := svc.GetItem(input)
	// Error when fetching
	if err != nil {
		// is it an AWS error?
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				log.Error(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				log.Error(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeRequestLimitExceeded:
				log.Error(dynamodb.ErrCodeRequestLimitExceeded, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				log.Error(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				log.Error(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			return "", errors.New(err.Error())
		}
	}

	// Grab the title from the result and return it
	// TODO:
	// optionally update ttl in DB -> frequent stuff gets cached longer

	if result.Item != nil {
		title := result.Item["title"].S
		return *title, nil
	}
	return "", errors.New("Cache miss")
}

// CacheAndReturn inserts a successfully found title to cache
func CacheAndReturn(query TitleQuery, title string, err error) (TitleQuery, error) {
	if err != nil {
		query.Title = ""
		return query, err
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1")},
	)

	if err != nil {
		log.Errorf("could not connect to AWS %v", err)
		return TitleQuery{}, err
	}

	// Create DynamoDB client
	svc := dynamodb.New(sess)
	// create a map for DD
	query.Title = title
	query.Added = time.Now().Unix()
	query.TTL = time.Now().Unix() + 86400 // 24 hours

	log.Infof("Storing TitleQuery: %v", query)

	av, err := dynamodbattribute.MarshalMap(query)
	if err != nil {
		log.Errorf("error marshaling to dynamodb: %v", err)
		return TitleQuery{}, err
	}

	// construct an input that DD canhandle
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("urls"),
	}
	// put item in DD
	_, err = svc.PutItem(input)

	return query, err
}
