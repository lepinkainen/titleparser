package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

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
	URL string `json:"url"`
}

type TitleResponse struct {
	Title string `json:"title"`
}

// HandleRequest is the function entry point
func HandleRequest(ctx context.Context, query TitleQuery) (TitleResponse, error) {
	url, err := FindTitle(query.URL)
	return TitleResponse{Title: url}, err
}

func main() {
	lambda.Start(HandleRequest)
}

func test() {

	urls := []string{"https://www.iltalehti.fi/kotimaa/a/bb594cd9-f66c-4bca-b626-a54848ea6ffb",
		"https://www.is.fi/taloussanomat/art-2000005940825.html?ref=rss",
		"https://www.is.fi/taloussanomat/oma-raha/art-2000005935228.html"}

	for _, url := range urls {
		fmt.Println(FindTitle(url))

	}
}
