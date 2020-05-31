package handler

import (
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/lepinkainen/titleparser/lambda"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var (
	// keyRegex attempts to find the cookie key we need to authenticate ourselves as a human :D
	keyRegex = regexp.MustCompile(`key=(.*?)\;`)
)

// handleYlilauta fetches Ylilauta titles by doing the cookie challenge correctly
func handleYlilauta(url string, authKey string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}

	// Set headers
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:49.0) Gecko/20100101 Firefox/49.0")
	req.Header.Set("Accept-Language", "*")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	// Create and Add cookie to request
	cookie := http.Cookie{Name: "key", Value: authKey}
	req.AddCookie(&cookie)

	// Set client timeout
	client := &http.Client{Timeout: time.Second * 10}

	// Send request
	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	defer res.Body.Close()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
		return "", errors.Wrap(err, "Could not load HTML")
	}

	// Authkey not set, fetch the challenge token
	if authKey == "" {
		log.Debugln("Fetching cookie for challenge")
		// Title wasn't found, do the cookie challenge
		cookiechallenge := doc.Find("script").Text()
		keymatches := keyRegex.FindStringSubmatch(cookiechallenge)
		if len(keymatches) != 0 {
			log.Debugf("Challenge cookie found: %s", keymatches[1])
			return handleYlilauta(url, keymatches[1])
		}
	}

	// We're in, get the title and get out
	s := doc.Find("title")
	if s != nil && s.Size() > 0 {
		title := s.First().Text()
		title = strings.TrimSuffix(title, " | Ylilauta")
		log.Debugf("Title found: %s", title)
		return title, nil
	}

	return "", errors.New("No title found from URL")
}

// Ylilauta handler, handles the cookie challenge
func Ylilauta(url string) (string, error) {
	return handleYlilauta(url, "")
}

func init() {
	lambda.RegisterHandler(".*?ylilauta.org/.*", Ylilauta)
}
