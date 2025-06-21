package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/lepinkainen/titleparser/lambda"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

var (
	apiURL = "https://www.googleapis.com/youtube/v3/videos"

	youtubeRegex1 = regexp.MustCompile(`https?://youtu.be/([^\?]+)([\?#]t=.*)?`)
	youtubeRegex2 = regexp.MustCompile(`https?://.*?youtube\.com/watch\?.*?v=([^&#]+)`)
)

type YoutubeReply struct {
	Etag  string `json:"etag"`
	Items []struct {
		ContentDetails struct {
			Caption       string `json:"caption"`
			ContentRating struct {
				YtRating string `json:"ytRating"`
			} `json:"contentRating"`
			Definition      string `json:"definition"`
			Dimension       string `json:"dimension"`
			Duration        string `json:"duration"`
			LicensedContent bool   `json:"licensedContent"`
			Projection      string `json:"projection"`
		} `json:"contentDetails"`
		Etag    string `json:"etag"`
		ID      string `json:"id"`
		Kind    string `json:"kind"`
		Snippet struct {
			CategoryID           string `json:"categoryId"`
			ChannelID            string `json:"channelId"`
			ChannelTitle         string `json:"channelTitle"`
			DefaultAudioLanguage string `json:"defaultAudioLanguage"`
			Description          string `json:"description"`
			LiveBroadcastContent string `json:"liveBroadcastContent"`
			Localized            struct {
				Description string `json:"description"`
				Title       string `json:"title"`
			} `json:"localized"`
			PublishedAt string   `json:"publishedAt"`
			Tags        []string `json:"tags"`
			Title       string   `json:"title"`
		} `json:"snippet"`
		Statistics struct {
			CommentCount  string `json:"commentCount"`
			DislikeCount  string `json:"dislikeCount"`
			FavoriteCount string `json:"favoriteCount"`
			LikeCount     string `json:"likeCount"`
			ViewCount     string `json:"viewCount"`
		} `json:"statistics"`
	} `json:"items"`
}

// OMDB handler
func Youtube(url string) (string, error) {
	youtubeKey := os.Getenv("YOUTUBE_KEY")
	if youtubeKey == "" {
		return "", errors.New("No API key set for Youtube")
	}

	match := youtubeRegex1.FindStringSubmatch(url)
	// Try again with second regex
	if len(match) == 0 {
		match = youtubeRegex2.FindStringSubmatch(url)
		if len(match) == 0 {
			return "", errors.New("Not a youtube URL")
		}
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.Fatal(err)
	}

	q := req.URL.Query()
	q.Add("id", match[1])
	q.Add("key", youtubeKey)
	q.Add("part", "snippet,contentDetails,statistics")
	q.Add("fields", "items(id,snippet,contentDetails,statistics)")
	req.URL.RawQuery = q.Encode()

	// Query the API
	res, err := client.Do(req)
	if err != nil {
		log.Error("Error querying Youtube API")
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		bytes, _ := io.ReadAll(res.Body)
		fmt.Println(string(bytes[:]))
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		return "", errors.Wrap(err, "HTTP error")
	}

	// error ignored on purpose
	bytes, _ := io.ReadAll(res.Body)

	// unmarshal JSON
	var reply YoutubeReply

	err = json.Unmarshal(bytes, &reply)
	if err != nil {
		return "", err
	}

	video := reply.Items[0]

	// The tag value is an ISO 8601 duration in the format PT#M#S
	duration := strings.ToLower(video.ContentDetails.Duration[2:])

	ageRestricted := ""
	if video.ContentDetails.ContentRating.YtRating == "ytAgeRestricted" {
		ageRestricted = " - age restricted"
	}

	publishedAt, err := time.Parse(time.RFC3339, video.Snippet.PublishedAt)
	if err != nil {
		log.Fatal(err)
	}
	agestr := humanize.RelTime(publishedAt, time.Now(), "ago", "from now")

	viewCount, _ := strconv.ParseInt(video.Statistics.ViewCount, 10, 64)
	views := HumanizeNumber(int(viewCount))

	// Ralph Breaks the Internet (2018)
	title := fmt.Sprintf("%s by %s [%s - %s views - %s%s]", video.Snippet.Title, video.Snippet.ChannelTitle, duration, views, agestr, ageRestricted)

	return title, nil
}

func init() {
	lambda.RegisterHandler(".*youtu.be.*", Youtube)
	lambda.RegisterHandler(".*youtube\\.com.*", Youtube)
}
