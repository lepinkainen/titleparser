package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"

	"github.com/lepinkainen/titleparser/lambda"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

var (
	// url pattern for OMDB searches
	// TODO: Grab the actual URL query param generation from youtube.go
	omdbURL = "http://www.omdbapi.com/?i=%s&apikey=%s"
	// figure out imdb id from url
	imdbRegex = regexp.MustCompile(`^https://www\.imdb\.com/title/(tt[\d]+)/?.*$`)
)

type omdbReply struct {
	Title   string `json:"Title"`
	Year    string `json:"Year"`
	Runtime string `json:"Runtime"`
	Genre   string `json:"Genre"`
	Ratings []struct {
		Source string `json:"Source"`
		Value  string `json:"Value"`
	} `json:"Ratings"`
	Metascore  string `json:"Metascore"`
	ImdbRating string `json:"imdbRating"`
	ImdbVotes  string `json:"imdbVotes"`
}

// OMDB handler
func OMDB(url string) (string, error) {
	omdbKey := os.Getenv("OMDB_KEY")
	if omdbKey == "" {
		return "", errors.New("No API key set for OMDB")
	}

	id := imdbRegex.FindStringSubmatch(url)
	if len(id) < 2 {
		return "", errors.New("No title ID found in URL")
	}

	// Request the HTML page.
	res, err := http.Get(fmt.Sprintf(omdbURL, id[1], omdbKey))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Warnf("Failed to close response body: %v", err)
		}
	}()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		return "", errors.Wrap(err, "HTTP error")
	}

	// error ignored on purpose
	bytes, _ := io.ReadAll(res.Body)

	// unmarshal JSON
	var reply omdbReply

	err = json.Unmarshal(bytes, &reply)
	if err != nil {
		return "", err
	}

	// Get possible scores from the Rating section
	var rtScore = "N/A"
	var imdbScore = "N/A"
	var metaScore = "N/A"

	for _, rating := range reply.Ratings {
		switch rating.Source {
		case "Rotten Tomatoes":
			rtScore = rating.Value
		case "Internet Movie Database":
			imdbScore = rating.Value
		case "Metacritic":
			metaScore = rating.Value
		}
	}

	// Ralph Breaks the Internet (2018)
	title := fmt.Sprintf("%s (%s) [IMDb %s] [RT %s] [Meta %s]",
		reply.Title, reply.Year, imdbScore, rtScore, metaScore)

	return title, nil
}

func init() {
	lambda.RegisterHandler(".*?imdb\\.com/title/tt.*", OMDB)
}
