package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gocolly/colly/extensions"

	"github.com/gocolly/colly"
)

var (
	ErrTitleNotFound = errors.New("No title found from URL")
)

func FindTitle(url string) (string, error) {

	var title string
	var ogTitle string

	// Instantiate default collector

	c := colly.NewCollector()
	c.IgnoreRobotsTxt = true
	extensions.RandomUserAgent(c)

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
		return ogTitle, nil
	} else if title != "" {
		return title, nil
	} else {
		return "", ErrTitleNotFound
	}
}

type TitleQuery struct {
	Url string `json:"url"`
}

type TitleResponse struct {
	Title string `json:"title"`
}

// HandleRequest is the function entry point
func HandleRequest(ctx context.Context, query TitleQuery) (TitleResponse, error) {
	url, err := FindTitle(query.Url)
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
