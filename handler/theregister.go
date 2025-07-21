package handler

import (
	"net/http"
	"time"

	"github.com/lepinkainen/titleparser/lambda"
	log "github.com/sirupsen/logrus"
)

var TheRegisterMatch = `.*\.theregister\.com.*|^https?://theregister\.com.*`

func TheRegister(url string) (string, error) {
	log.Infof("Using The Register handler for %s", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Errorf("Error creating request for %s: %v", url, err)
		return "", err
	}

	// Set headers to avoid 403 Forbidden from The Register's bot detection
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Accept-Language", AcceptLanguage)
	req.Header.Set("Accept", Accept)

	// Set client timeout (10 seconds, same as other handlers)
	client := &http.Client{Timeout: time.Second * 10}

	// Send request
	res, err := client.Do(req)
	if err != nil {
		log.Errorf("Error sending request to %s: %v", url, err)
		return "", err
	}
	defer res.Body.Close()

	log.Infof("The Register response status: %d for %s", res.StatusCode, url)

	// Use the default parsing logic from lambda.DefaultHandler
	// We can't call DefaultHandler directly since it does its own HTTP request,
	// so we need to use the same parsing logic here
	return lambda.ParseHTMLFromResponse(res, url)
}

func init() {
	lambda.RegisterHandler(TheRegisterMatch, TheRegister)
}
