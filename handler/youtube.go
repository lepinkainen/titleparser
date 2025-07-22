package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
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
	videoAPIURL   = "https://www.googleapis.com/youtube/v3/videos"
	channelAPIURL = "https://www.googleapis.com/youtube/v3/channels"

	// Video URL patterns
	youtubeRegex1 = regexp.MustCompile(`https?://youtu.be/([^\?]+)([\?#]t=.*)?`)
	youtubeRegex2 = regexp.MustCompile(`https?://.*?youtube\.com/watch\?.*?v=([^&#]+)`)

	// Channel URL patterns
	channelHandleRegex = regexp.MustCompile(`https?://.*?youtube\.com/@([^/?#]+)`)
	channelCustomRegex = regexp.MustCompile(`https?://.*?youtube\.com/c/([^/?#]+)`)
	channelIDRegex     = regexp.MustCompile(`https?://.*?youtube\.com/channel/([^/?#]+)`)
	channelUserRegex   = regexp.MustCompile(`https?://.*?youtube\.com/user/([^/?#]+)`)
	channelDirectRegex = regexp.MustCompile(`https?://.*?youtube\.com/([^/?#]+)$`)
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

type YoutubeChannelReply struct {
	Etag  string `json:"etag"`
	Items []struct {
		Etag    string `json:"etag"`
		ID      string `json:"id"`
		Kind    string `json:"kind"`
		Snippet struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			CustomURL   string `json:"customUrl"`
			PublishedAt string `json:"publishedAt"`
			Localized   struct {
				Title       string `json:"title"`
				Description string `json:"description"`
			} `json:"localized"`
		} `json:"snippet"`
		Statistics struct {
			ViewCount             string `json:"viewCount"`
			SubscriberCount       string `json:"subscriberCount"`
			HiddenSubscriberCount bool   `json:"hiddenSubscriberCount"`
			VideoCount            string `json:"videoCount"`
		} `json:"statistics"`
	} `json:"items"`
}

func Youtube(url string) (string, error) {
	youtubeKey := os.Getenv("YOUTUBE_KEY")
	if youtubeKey == "" {
		return "", errors.New("No API key set for Youtube")
	}

	// Check if this is a video URL
	if videoID := ExtractVideoID(url); videoID != "" {
		return handleVideoURL(videoID, youtubeKey)
	}

	// Check if this is a channel URL
	if channelID, paramType := ExtractChannelInfo(url); channelID != "" {
		return handleChannelURL(channelID, paramType, youtubeKey)
	}

	return "", errors.New("Not a valid YouTube URL")
}

func ExtractVideoID(url string) string {
	if match := youtubeRegex1.FindStringSubmatch(url); len(match) > 1 {
		return match[1]
	}
	if match := youtubeRegex2.FindStringSubmatch(url); len(match) > 1 {
		return match[1]
	}
	return ""
}

func ExtractChannelInfo(url string) (string, string) {
	if match := channelHandleRegex.FindStringSubmatch(url); len(match) > 1 {
		return match[1], "forHandle"
	}
	if match := channelCustomRegex.FindStringSubmatch(url); len(match) > 1 {
		return match[1], "forUsername"
	}
	if match := channelIDRegex.FindStringSubmatch(url); len(match) > 1 {
		return match[1], "id"
	}
	if match := channelUserRegex.FindStringSubmatch(url); len(match) > 1 {
		return match[1], "forUsername"
	}
	if match := channelDirectRegex.FindStringSubmatch(url); len(match) > 1 {
		// Direct channel URLs could be either custom names or handles
		return match[1], "forUsername"
	}
	return "", ""
}

func handleVideoURL(videoID, apiKey string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", videoAPIURL, nil)
	if err != nil {
		return "", errors.Wrap(err, "Failed to create video API request")
	}

	q := req.URL.Query()
	q.Add("id", videoID)
	q.Add("key", apiKey)
	q.Add("part", "snippet,contentDetails,statistics")
	q.Add("fields", "items(id,snippet,contentDetails,statistics)")
	req.URL.RawQuery = q.Encode()

	res, err := client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "Error querying YouTube video API")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		bytes, _ := io.ReadAll(res.Body)
		log.Errorf("YouTube API error: %d %s - %s", res.StatusCode, res.Status, string(bytes))
		return "", errors.New("YouTube API returned error status")
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", errors.Wrap(err, "Failed to read video API response")
	}

	var reply YoutubeReply
	if err := json.Unmarshal(bytes, &reply); err != nil {
		return "", errors.Wrap(err, "Failed to parse video API response")
	}

	if len(reply.Items) == 0 {
		return "", errors.New("Video not found")
	}

	video := reply.Items[0]

	duration := strings.ToLower(video.ContentDetails.Duration[2:])

	ageRestricted := ""
	if video.ContentDetails.ContentRating.YtRating == "ytAgeRestricted" {
		ageRestricted = " - age restricted"
	}

	publishedAt, err := time.Parse(time.RFC3339, video.Snippet.PublishedAt)
	if err != nil {
		log.Warnf("Failed to parse video publish date: %v", err)
		return fmt.Sprintf("%s by %s", video.Snippet.Title, video.Snippet.ChannelTitle), nil
	}
	agestr := humanize.RelTime(publishedAt, time.Now(), "ago", "from now")

	viewCount, _ := strconv.ParseInt(video.Statistics.ViewCount, 10, 64)
	var views string
	if viewCount >= math.MinInt && viewCount <= math.MaxInt {
		views = HumanizeNumber(int(viewCount))
	} else {
		log.Warnf("View count exceeds int range: %d", viewCount)
		views = HumanizeNumber(math.MaxInt)
	}

	title := fmt.Sprintf("%s by %s [%s - %s views - %s%s]",
		video.Snippet.Title, video.Snippet.ChannelTitle, duration, views, agestr, ageRestricted)

	return title, nil
}

func handleChannelURL(channelID, paramType, apiKey string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", channelAPIURL, nil)
	if err != nil {
		return "", errors.Wrap(err, "Failed to create channel API request")
	}

	q := req.URL.Query()
	q.Add(paramType, channelID)
	q.Add("key", apiKey)
	q.Add("part", "snippet,statistics")
	q.Add("fields", "items(id,snippet,statistics)")
	req.URL.RawQuery = q.Encode()

	res, err := client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "Error querying YouTube channel API")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		bytes, _ := io.ReadAll(res.Body)
		log.Errorf("YouTube channel API error: %d %s - %s", res.StatusCode, res.Status, string(bytes))
		return "", errors.New("YouTube channel API returned error status")
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", errors.Wrap(err, "Failed to read channel API response")
	}

	var reply YoutubeChannelReply
	if err := json.Unmarshal(bytes, &reply); err != nil {
		return "", errors.Wrap(err, "Failed to parse channel API response")
	}

	if len(reply.Items) == 0 {
		return "", errors.New("Channel not found")
	}

	channel := reply.Items[0]

	subscriberCount, _ := strconv.ParseInt(channel.Statistics.SubscriberCount, 10, 64)
	var subscribers string
	if subscriberCount >= math.MinInt && subscriberCount <= math.MaxInt {
		subscribers = HumanizeNumber(int(subscriberCount))
	} else {
		log.Warnf("Subscriber count exceeds int range: %d", subscriberCount)
		subscribers = HumanizeNumber(math.MaxInt)
	}

	videoCount, _ := strconv.ParseInt(channel.Statistics.VideoCount, 10, 64)
	var videos string
	if videoCount >= math.MinInt && videoCount <= math.MaxInt {
		videos = HumanizeNumber(int(videoCount))
	} else {
		log.Warnf("Video count exceeds int range: %d", videoCount)
		videos = HumanizeNumber(math.MaxInt)
	}

	publishedAt, err := time.Parse(time.RFC3339, channel.Snippet.PublishedAt)
	if err != nil {
		log.Warnf("Failed to parse channel publish date: %v", err)
		return fmt.Sprintf("%s [Channel - %s subscribers]", channel.Snippet.Title, subscribers), nil
	}
	agestr := humanize.RelTime(publishedAt, time.Now(), "ago", "from now")

	title := fmt.Sprintf("%s [Channel - %s subscribers - %s videos - created %s]",
		channel.Snippet.Title, subscribers, videos, agestr)

	return title, nil
}

func init() {
	lambda.RegisterHandler(".*youtu.be.*", Youtube)
	lambda.RegisterHandler(".*youtube\\.com.*", Youtube)
}
