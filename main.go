package main

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"titleparser/handler"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gocolly/colly/extensions"

	"github.com/gocolly/colly"
)

var (
	// ErrTitleNotFound is returned when the target resource doesn't have a title
	ErrTitleNotFound = errors.New("No title found from URL")
	// ErrNotHTML is returned when the source url is not of type text/html
	ErrNotHTML = errors.New("Source url is not HTML")
)

func collyError(r *colly.Response, err error) {
	fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
}

// FindTitle returns the title or opengraph title of the given url
func FindTitle(url string) (string, error) {

	var title string
	var ogTitle string
	var err error

	// Instantiate default collector

	c := colly.NewCollector()

	c.IgnoreRobotsTxt = true
	c.MaxBodySize = 1024 * 1024 // 1MB maximum

	extensions.RandomUserAgent(c)

	c.OnError(collyError)

	// Before making a request print "Visiting ..."
	c.OnResponse(func(r *colly.Response) {
		contentType := r.Headers.Get("Content-Type")
		if !strings.HasPrefix(contentType, "text/html") {
			fmt.Printf("Invalid content type: %s\n", contentType)
			err = ErrNotHTML
		}
	})

	// Find the regular title
	// <title>
	c.OnHTML("title", func(e *colly.HTMLElement) {
		title = e.Text
	})

	// Opengraph title
	// <meta property="og:title" content="Title" />
	c.OnHTML("meta", func(e *colly.HTMLElement) {
		if !(e.Attr("property") == "og:title") {
			return
		}
		ogTitle = e.Attr("content")
	})

	c.Visit(url)

	// prefer og:title, since it tends to have less crap in it
	if ogTitle != "" {
		return ogTitle, err
	} else if title != "" {
		return title, err
	} else {
		return "", ErrTitleNotFound
	}
}

type TitleQuery struct {
	User    string `json:"user"`
	Channel string `json:"channel"`
	URL     string `json:"url"`
}

type TitleResponse struct {
	Title string `json:"title"`
}

// CheckCache will return a non-empty string if the URL given is in the cache
func CheckCache(query TitleQuery) string {
	// TODO:
	// connect to dynamodb
	// attempt to fetch url with query.URL
	// return title
	// optionally update ttl in DB
	return ""
}

// CacheAndReturn inserts a successfully found title to cache
func CacheAndReturn(query TitleQuery, title string, err error) (TitleResponse, error) {
	// insert url to cache

	return TitleResponse{Title: title}, err
}

// HandleRequest is the function entry point
func HandleRequest(ctx context.Context, query TitleQuery) (TitleResponse, error) {

	// if query is cached, return from cache instead of fetching
	if title := CheckCache(query); title != "" {
		return CacheAndReturn(query, title, nil)
	}

	// https://golang.org/pkg/path/filepath/#Match
	handlerFunctions := make(map[string]func(string) (string, error))

	handlerFunctions[".*?areena.yle.fi/.*"] = handler.YleAreena
	handlerFunctions[".*?apina.biz.*"] = handler.ApinaBiz

	for pattern, handler := range handlerFunctions {
		match, err := regexp.MatchString(pattern, query.URL)

		// error in matching
		if err != nil {
			fmt.Printf("Error matching with pattern %s", pattern)
			return TitleResponse{Title: ""}, err
		}

		// no error and match, run function to get actual title
		if err == nil && match {
			title, err := handler(query.URL)
			return CacheAndReturn(query, title, err)
		}
	}

	url, err := FindTitle(query.URL)
	return TitleResponse{Title: url}, err
}

func main() {
	lambda.Start(HandleRequest)
}
