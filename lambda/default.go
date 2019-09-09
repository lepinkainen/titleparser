package lambda

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var (
	// ErrTitleNotFound is returned when the target resource doesn't have a title
	ErrTitleNotFound = errors.New("No title found from URL")
	// ErrNotHTML is returned when the source url is not of type text/html
	ErrNotHTML = errors.New("Source url is not HTML")

	// TitleMax is the maximum length for a title
	TitleMax = 200
)

// DefaultHandler is the fallback for sites that don't have a special handler
// TODO: Split to two parts: 1) fetch url 2) parse title from html
// 		 Tests for both parts
func DefaultHandler(url string) (string, error) {
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	// Not html, don't bother parsing
	contentType := res.Header.Get("content-type")
	if !strings.HasPrefix(contentType, "text/html") {
		return "", ErrNotHTML
	}

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		return "", errors.Wrap(err, "HTTP error")
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
		return "", errors.Wrap(err, "Could not load HTML")
	}

	// primarily we want to use og:title
	s := doc.Find(`meta[property="og:title"]`)
	if s != nil && s.Size() > 0 {
		title, _ := s.Attr("content")
		return sanitize(title), nil
	}

	// Bleh, just a boring old title then
	s = doc.Find("title")
	if s != nil && s.Size() > 0 {
		// Just grab the first one, some pages (ab)use the title element
		title := s.First().Text()
		return sanitize(title), nil
	}

	// No title, report it
	return "", ErrTitleNotFound
}

// sanitize the url by removing everything superfluous
func sanitize(title string) string {
	// remove newlines
	title = strings.ReplaceAll(title, "\n", "")
	title = strings.ReplaceAll(title, "\r", "")
	// remove leading and trailing whitespace
	title = strings.TrimSpace(title)

	// max size 200 characters. It's a title, not a goddamn novel
	end := len(title)
	if end > TitleMax {
		end = TitleMax
		title = fmt.Sprintf("%s...", title[:end])
	}

	return title
}
