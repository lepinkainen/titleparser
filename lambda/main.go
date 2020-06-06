package lambda

import (
	"context"
	"os"
	"regexp"
	"time"

	log "github.com/sirupsen/logrus"
	//"github.com/lepinkainen/titleparser/handler"
)

var (
	handlerFunctions = make(map[string]func(string) (string, error))
)

// TitleQuery received via HTTP(s)
type TitleQuery struct {
	Added   int64  `json:"timestamp"`
	User    string `json:"user"`
	Channel string `json:"channel"`
	URL     string `json:"url"`
	Title   string `json:"title"`
	TTL     int64  `json:"ttl"` // TTL is used to expire the item in DynamoDB automatically
}

type handlerFunc func(string) (string, error)

// RegisterHandler adds the given url parser and pattern to the map of handlers
func RegisterHandler(pattern string, function handlerFunc) {
	handlerFunctions[pattern] = function
}

// func registerParser(pattern string, function handlerFunc) {
// 	handlerFunctions[pattern] = function
// }

// HandleRequest is the function entry point
func HandleRequest(ctx context.Context, query TitleQuery) (TitleQuery, error) {

	// open session to dynamoDB

	log.Infof("Handling %v", query)

	// If we are running locally, don't use dynamodb as a cache
	// TODO: Possibly add an in-memory DB or sqlite for local mode caching?
	var runmode = os.Getenv("RUNMODE")
	if runmode != "local" {
		// if query is cached, return from cache instead of fetching
		if title, err := CheckCache(query); err == nil {
			return CacheAndReturn(query, title, nil)
		}
	}

	for pattern, handler := range handlerFunctions {
		match, err := regexp.MatchString(pattern, query.URL)

		// error in matching, log and continue
		if err != nil {
			log.Errorf("Error matching with pattern %s: %v", pattern, err)
		}

		// no error and match, run function to get actual title and return
		if err == nil && match {
			log.Infof("Handler match found for %s\n", query.URL)
			title, err := handler(query.URL)
			if runmode != "local" {
				return CacheAndReturn(query, title, err)
			}
			log.Infoln("Local mode, not caching result")

			query.Title = title
			query.Added = time.Now().Unix()
			query.TTL = time.Now().Unix() + 86400 // 24 hours

			return query, err
		}
	}

	log.Infof("No handler found for %s, falling back to default", query.URL)

	// custom parsers didn't match, use the default parser
	title, err := DefaultHandler(query.URL)
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
