package handler

import (
	"net/http"
	"time"

	"github.com/lepinkainen/titleparser/lambda"
	log "github.com/sirupsen/logrus"
)

// ThreadsMatch matches Threads share links, post URLs and short /t/ links on
// both the threads.com and legacy threads.net domains.
var ThreadsMatch = `.*(//|\.)threads\.(net|com)/(share/|t/|@[^/]+/post/).*`

// ThreadsUserAgent is a social-crawler User-Agent. Threads serves the client
// JS shell (with only "<title>Threads</title>" and no OpenGraph tags) to
// regular browser User-Agents, but returns full og:title/og:image metadata to
// recognised social crawlers. We impersonate one to get a real title.
const ThreadsUserAgent = "Mozilla/5.0 (compatible; Twitterbot/1.0)"

// Threads extracts the title for a Threads URL by requesting the page as a
// social crawler so that the server includes OpenGraph metadata.
func Threads(url string) (string, error) {
	log.Infof("Using Threads handler for %s", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Errorf("Error creating request for %s: %v", url, err)
		return "", err
	}

	// Crawler User-Agent is required, otherwise Threads returns the JS shell
	// without OpenGraph tags.
	req.Header.Set("User-Agent", ThreadsUserAgent)
	req.Header.Set("Accept-Language", AcceptLanguage)
	req.Header.Set("Accept", Accept)

	client := &http.Client{Timeout: time.Second * 10}

	res, err := client.Do(req)
	if err != nil {
		log.Errorf("Error sending request to %s: %v", url, err)
		return "", err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Warnf("Failed to close response body: %v", err)
		}
	}()

	log.Infof("Threads response status: %d for %s", res.StatusCode, url)

	// Reuse the shared OpenGraph/title parsing logic.
	return lambda.ParseHTMLFromResponse(res, url)
}

func init() {
	lambda.RegisterHandler(ThreadsMatch, Threads)
}
