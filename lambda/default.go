package lambda

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
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

// DefaultHandler returns the title or opengraph title of the given url
func DefaultHandler(url string) (string, error) {

	var title string
	var ogTitle string
	var err error

	// Instantiate default collector

	c := colly.NewCollector()

	c.IgnoreRobotsTxt = true
	c.MaxBodySize = 1024 * 1024 // 1MB maximum

	extensions.RandomUserAgent(c)

	c.OnError(collyError)

	// Check that the URL is actually something parseable
	c.OnResponse(func(r *colly.Response) {
		contentType := r.Headers.Get("Content-Type")
		if !strings.HasPrefix(contentType, "text/html") {
			fmt.Printf("Invalid content type: %s\n", contentType)
			err = ErrNotHTML
		}
	})

	// Find the <title> tag
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
