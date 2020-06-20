package handler

import (
	"context"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dustin/go-humanize"
	"github.com/lepinkainen/titleparser/lambda"
	"golang.org/x/oauth2/clientcredentials"
)

var (
	// Single tweet page
	twitterRegex = regexp.MustCompile(`.*?twitter\.com/.*?/status/(\d+)`)

	// Custom magnitudes for humanize, shorter and more concise
	shortMagnitudes = []humanize.RelTimeMagnitude{
		{D: time.Second, Format: "now", DivBy: time.Second},
		{D: 2 * time.Second, Format: "1s", DivBy: 1},
		{D: time.Minute, Format: "%ds", DivBy: time.Second},
		{D: 2 * time.Minute, Format: "1m", DivBy: 1},
		{D: time.Hour, Format: "%dm", DivBy: time.Minute},
		{D: 2 * time.Hour, Format: "1h", DivBy: 1},
		{D: humanize.Day, Format: "%dh", DivBy: time.Hour},
		{D: 2 * humanize.Day, Format: "1d", DivBy: 1},
		{D: humanize.Week, Format: "%dd", DivBy: humanize.Day},
		{D: 2 * humanize.Week, Format: "1w", DivBy: 1},
		{D: humanize.Month, Format: "%dw", DivBy: humanize.Week},
		{D: 2 * humanize.Month, Format: "1m", DivBy: 1},
		{D: humanize.Year, Format: "%dm", DivBy: humanize.Month},
		{D: 18 * humanize.Month, Format: "1y", DivBy: 1},
		{D: 2 * humanize.Year, Format: "2y", DivBy: 1},
		{D: humanize.LongTime, Format: "%dy", DivBy: humanize.Year},
		{D: math.MaxInt64, Format: "a long long time ago in a galaxy far far away", DivBy: 1},
	}
)

// Twitter single tweet urls
func Twitter(url string) (string, error) {
	groups := twitterRegex.FindStringSubmatch(url)
	tweetID, err := strconv.ParseInt(groups[1], 10, 64)
	if err != nil {
		log.Error("Could not parse tweet id: ", err)
		return "", err
	}

	// oauth2 configures a client that uses app credentials to keep a fresh token
	config := &clientcredentials.Config{
		ClientID:     os.Getenv("TWITTER_CLIENTID"),
		ClientSecret: os.Getenv("TWITTER_CLIENTSECRET"),
		TokenURL:     "https://api.twitter.com/oauth2/token",
	}
	// http.Client will automatically authorize Requests
	httpClient := config.Client(context.Background())

	// Twitter client
	client := twitter.NewClient(httpClient)

	tweet, _, err := client.Statuses.Show(tweetID, &twitter.StatusShowParams{TweetMode: "extended"})
	if err != nil {
		log.Errorf("Error fetching tweet: %s", err)
		return "", err
	}

	var verifiedMark = ""
	if tweet.User.Verified {
		verifiedMark = "✔"
	}

	// TODO: Handle tweet age
	createTime, err := tweet.CreatedAtTime()
	if err != nil {
		log.Errorf("Could not get tweet creation time: %s", err)
		return "", err
	}

	// TODO: 2 weeks ago -> 2w would require a custom magnitudes -struct
	ago := humanize.CustomRelTime(createTime, time.Now(), "", "", shortMagnitudes)

	//ago := humanize.Time(createTime)

	tweetString := fmt.Sprintf("%s (%s@%s) %s: %s [♻ %d ♥ %d]", tweet.User.Name, verifiedMark, tweet.User.ScreenName, ago, tweet.FullText, tweet.RetweetCount, tweet.FavoriteCount)

	return tweetString, nil
}

// Register the handler function with corresponding regex
func init() {
	lambda.RegisterHandler(`.*?twitter\.com/.*?/status/.*?`, Twitter)
}
