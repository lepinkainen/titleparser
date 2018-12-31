package main

import (
	"context"
	"os"
	"regexp"

	"titleparser/handler"

	"github.com/aws/aws-lambda-go/lambda"

	log "github.com/sirupsen/logrus"
)

var (
	handlerFunctions = make(map[string]func(string) (string, error))
)

type TitleQuery struct {
	Added   int64  `json:"timestamp"`
	User    string `json:"user"`
	Channel string `json:"channel"`
	URL     string `json:"url"`
	Title   string `json:"title"`
	TTL     int64  `json:"ttl"` // TTL is used to expire the item in DynamoDB automatically
}

type handlerFunc func(string) (string, error)

// RegisterParser adds the given url parser and pattern to the map of handlers
func RegisterParser(pattern string, function handlerFunc) {
	handlerFunctions[pattern] = function
}

// HandleRequest is the function entry point
func HandleRequest(ctx context.Context, query TitleQuery) (TitleQuery, error) {

	// open session to dynamoDB

	log.Infof("Handling %v", query)

	// if query is cached, return from cache instead of fetching
	if title, err := CheckCache(query); err == nil {
		return CacheAndReturn(query, title, nil)
	}

	// register custom parsers
	RegisterParser(".*?areena.yle.fi/.*", handler.YleAreena)
	RegisterParser(".*?apina.biz.*", handler.ApinaBiz)
	//RegisterParser(".*", handler.DefaultHandler)

	for pattern, handler := range handlerFunctions {
		match, err := regexp.MatchString(pattern, query.URL)

		// error in matching, log and continue
		if err != nil {
			log.Errorf("Error matching with pattern %s: %v", pattern, err)
		}

		// no error and match, run function to get actual title and return
		if err == nil && match {
			title, err := handler(query.URL)
			return CacheAndReturn(query, title, err)
		}
	}

	// custom parsers didn't match, use the default parser
	title, err := handler.DefaultHandler(query.URL)
	return CacheAndReturn(query, title, err)
}

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

func main() {
	lambda.Start(HandleRequest)
}
