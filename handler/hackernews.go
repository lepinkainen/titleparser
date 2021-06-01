package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"github.com/lepinkainen/titleparser/lambda"
	log "github.com/sirupsen/logrus"
)

// Fetch titles from hackernews using their API
// https://github.com/HackerNews/API

var (
	hnRegex  = regexp.MustCompile(`news\.ycombinator\.com\/item\?id=(\d+)`)
	hnAPIURL = "https://hacker-news.firebaseio.com/v0/item/%s.json"
)

// HNAPIResponse is the response json from Firebase
type HNAPIResponse struct {
	By          string `json:"by"`
	Descendants int    `json:"descendants"`
	ID          int    `json:"id"`
	Kids        []int  `json:"kids"`
	Score       int    `json:"score"`
	Time        int    `json:"time"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	URL         string `json:"url"`
}

// HackerNews titles using the API
func HackerNews(url string) (string, error) {
	storyID := hnRegex.FindStringSubmatch(url)

	if storyID == nil || len(storyID) < 2 {
		return "", nil
	}

	url = fmt.Sprintf(hnAPIURL, storyID[1])

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("Error reading request. ", err)
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
		log.Fatal("Error reading response. ", err)
	}
	defer res.Body.Close()

	var apiResponse HNAPIResponse

	body, _ := ioutil.ReadAll(res.Body)

	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		_ = fmt.Errorf("unable to unmarshal JSON response: %#v", err)
		return "", err
	}

	// TODO: Maybe handle other response types?
	// Comment link could dig up until it finds the top and pick title from there.

	// These are the ones, that have title and score
	if apiResponse.Type == "story" ||
		apiResponse.Type == "poll" ||
		apiResponse.Type == "job" {
		return fmt.Sprintf("%s by %s [%d points]", apiResponse.Title, apiResponse.By, apiResponse.Score), nil
	}

	return "", nil
}

// Register the handler function with corresponding regex
func init() {
	lambda.RegisterHandler(".*?news.ycombinator.com.*", HackerNews)
}
