package handler

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/dustin/go-humanize"
	"github.com/lepinkainen/titleparser/lambda"
	log "github.com/sirupsen/logrus"
)

// MastodonMatch matches URLs with the Mastodon post pattern: /@username/numeric_id
// This pattern works across all Mastodon instances
var MastodonMatch = ".*/@[^/]+/[0-9]+"

// Mastodon extracts the title from a Mastodon post URL
func Mastodon(url string) (string, error) {
	// Create a request with proper headers
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error("Error creating request: ", err)
		return "", err
	}

	// Set headers
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Accept-Language", AcceptLanguage)
	req.Header.Set("Accept", Accept)

	// Set client timeout
	client := &http.Client{Timeout: time.Second * 10}

	// Send request
	res, err := client.Do(req)
	if err != nil {
		log.Error("Error sending request: ", err)
		return "", err
	}
	defer res.Body.Close()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Error("Error parsing HTML: ", err)
		return "", err
	}

	// Extract relevant information
	var username, content, timestamp string
	var postTime time.Time

	// Extract username from URL using regex
	usernameRegex := regexp.MustCompile(`.*/@([^/]+)/[0-9]+`)
	matches := usernameRegex.FindStringSubmatch(url)
	if len(matches) > 1 {
		username = matches[1]
	}

	// Try to find post content (most Mastodon instances use article.status-content)
	content = doc.Find("article .status-content").Text()
	content = strings.TrimSpace(content)

	// Try to get timestamp
	timestampElement := doc.Find("article time")
	if timestampElement.Length() > 0 {
		if datetime, exists := timestampElement.Attr("datetime"); exists {
			// Parse the timestamp
			parsed, err := time.Parse(time.RFC3339, datetime)
			if err == nil {
				postTime = parsed
				timestamp = humanize.RelTime(postTime, time.Now(), "ago", "")
			}
		}
	}

	// If we couldn't get the content properly, fall back to og:description
	if content == "" {
		doc.Find(`meta[property="og:description"]`).Each(func(_ int, s *goquery.Selection) {
			if c, exists := s.Attr("content"); exists && c != "" {
				content = c
			}
		})
	}

	// Format the title
	var title string
	if content != "" {
		// Truncate content if too long
		if len(content) > 100 {
			content = content[:97] + "..."
		}

		if timestamp != "" {
			title = fmt.Sprintf("@%s: %s [%s]", username, content, timestamp)
		} else {
			title = fmt.Sprintf("@%s: %s", username, content)
		}
	} else {
		// Fallback to og:title if we couldn't extract content
		title = doc.Find(`meta[property="og:title"]`).AttrOr("content", "")
		if title == "" {
			title = doc.Find("title").Text()
		}
	}

	return title, nil
}

func init() {
	lambda.RegisterHandler(MastodonMatch, Mastodon)
}
