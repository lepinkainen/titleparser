package handler

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/lepinkainen/titleparser/lambda"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

// Verkkokauppa handler
func Verkkokauppa(url string) (string, error) {
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
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
		return title, nil
	}
	return "", errors.New("og:title not found")
}

func init() {
	lambda.RegisterHandler(`.*?verkkokauppa\.com/.*?/product/.*?`, Verkkokauppa)
}
