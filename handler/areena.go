package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dustin/go-humanize"

	"github.com/PuerkitoBio/goquery"
	"github.com/lepinkainen/titleparser/lambda"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

/*
   <meta property="og:title" content="Muiden elämä">
   <meta property="og:video:duration" content="7920">
   <meta property="og:video:release_date" content="2021-02-10T06:00:00.000+02:00">
*/

// YleAreena handler TBD
func YleAreena(url string) (string, error) {

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
		return "", errors.Wrap(err, "Could not load HTML")
	}

	var title, duration, release_date string

	// primarily we want to use og:title
	s := doc.Find(`meta[property="og:title"]`)
	if s != nil && s.Size() > 0 {
		title, _ = s.Attr("content")
	}

	s = doc.Find(`meta[property="og:video:duration"]`)
	if s != nil && s.Size() > 0 {
		s_content, _ := s.Attr("content")
		duration_time, _ := time.ParseDuration(fmt.Sprintf("%ss", s_content))
		duration = duration_time.String()
	}
	s = doc.Find(`meta[property="og:video:release_date"]`)
	if s != nil && s.Size() > 0 {
		s_content, _ := s.Attr("content")
		// 									  2018-12-25T06:00:00.000+02:00
		release_date_time, err := time.Parse("2006-01-02T15:04:05.000Z07:00", s_content)
		if err != nil {
			log.Error("Error parsing release date: ", err)
		}
		release_date = humanize.RelTime(release_date_time, time.Now(), "ago", "")

	}

	if duration == "" || release_date == "" {
		return title, nil
	} else {
		return fmt.Sprintf("%s [Duration: %s Released: %s]", title, duration, release_date), nil
	}
}

func init() {
	lambda.RegisterHandler(".*?areena.yle.fi/.*", YleAreena)
}
