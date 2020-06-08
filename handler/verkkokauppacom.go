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

	selection := doc.Find("h1")
	title := selection.Contents().Text()

	return title, nil
}

func init() {
	lambda.RegisterHandler("https://www.verkkokauppa.com/.*?/product/*.?", Verkkokauppa)
}
