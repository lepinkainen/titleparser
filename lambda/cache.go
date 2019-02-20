package lambda

import (
	"errors"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	log "github.com/sirupsen/logrus"
)

// CheckCache will return a non-empty string if the URL given is in the cache
func CheckCache(query TitleQuery) (string, error) {
	// TODO:
	// connect to dynamodb
	// attempt to fetch url with query.URL
	// return title
	// optionally update ttl in DB
	return "", errors.New("cache miss")
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
