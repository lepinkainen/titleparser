package lambda

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/lepinkainen/titleparser/common"
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

	// RemoveWhitespaceRegex marks all tokens with more than one whitespace
	RemoveWhitespaceRegex = regexp.MustCompile(`[\s]{2,}`)
)

// DefaultHandler is the fallback for sites that don't have a special handler
// TODO: Split to two parts: 1) fetch url 2) parse title from html
//
//	Tests for both parts
func DefaultHandler(url string) (string, error) {
	// Create request with proper browser headers to avoid User-Agent blocking
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Set headers to avoid 403 Forbidden from sites that block Go client
	req.Header.Set("User-Agent", common.UserAgent)
	req.Header.Set("Accept-Language", common.AcceptLanguage)
	req.Header.Set("Accept", common.Accept)

	// Set client timeout (10 seconds, consistent with other handlers)
	client := &http.Client{Timeout: time.Second * 10}

	// Send request
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Warnf("Failed to close response body: %v", err)
		}
	}()

	// Not html, don't bother parsing
	contentType := res.Header.Get("content-type")
	if !strings.HasPrefix(contentType, "text/html") {
		return "", ErrNotHTML
	}

	if res.StatusCode != 200 {
		switch res.StatusCode {
		case 403:
			return "", errors.New("403 Forbidden")
		case 404:
			return "", errors.New("404 Not Found")
		case 405:
			return "", errors.New("405 Method Not Allowed")
		case 429:
			return "", errors.New("429 Too Many Requests")
		case 500:
			return "", errors.New("500 Internal Server Error")
		case 502:
			return "", errors.New("502 Bad Gateway")
		default:
			log.Fatalf("unhandled status code: %d (%s)", res.StatusCode, res.Status)
			return "", errors.Wrap(err, "HTTP error")
		}
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
	// remove extra whitespace in the middle of the title
	// some crappy CMSes leave it all over the place
	title = RemoveWhitespaceRegex.ReplaceAllLiteralString(title, " ")
	// remove leading and trailing whitespace
	title = strings.TrimSpace(title)
	// remove newlines
	title = strings.ReplaceAll(title, "\n", "")
	title = strings.ReplaceAll(title, "\r", "")

	// max size 200 characters. It's a title, not a goddamn novel
	end := len(title)
	if end > TitleMax {
		end = TitleMax
		title = fmt.Sprintf("%s...", title[:end])
	}

	return title
}

// ParseHTMLFromResponse extracts title from an HTTP response
// This is used by custom handlers that need to do their own HTTP requests
func ParseHTMLFromResponse(res *http.Response, url string) (string, error) {
	// Not html, don't bother parsing
	contentType := res.Header.Get("content-type")
	if !strings.HasPrefix(contentType, "text/html") {
		return "", ErrNotHTML
	}

	if res.StatusCode != 200 {
		switch res.StatusCode {
		case 403:
			return "", errors.New("403 Forbidden")
		case 404:
			return "", errors.New("404 Not Found")
		case 405:
			return "", errors.New("405 Method Not Allowed")
		case 429:
			return "", errors.New("429 Too Many Requests")
		case 500:
			return "", errors.New("500 Internal Server Error")
		case 502:
			return "", errors.New("502 Bad Gateway")
		default:
			log.Errorf("unhandled status code: %d (%s) for URL: %s", res.StatusCode, res.Status, url)
			return "", fmt.Errorf("HTTP error: %d %s", res.StatusCode, res.Status)
		}
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Errorf("Could not load HTML from %s: %v", url, err)
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
