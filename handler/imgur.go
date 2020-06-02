package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/lepinkainen/titleparser/lambda"

	log "github.com/sirupsen/logrus"
)

var (
	// imgur gallery
	galleryRegex = regexp.MustCompile(`.*?imgur.com/gallery/(.*)`)
	// imgur album
	albumRegex = regexp.MustCompile(`.*?imgur.com/a/(.*)`)
)

// ImgurResponse is the imgur generic API response for all gallery queries
type ImgurResponse struct {
	Data struct {
		ID            string      `json:"id"`
		Title         string      `json:"title"`
		Description   interface{} `json:"description"`
		Datetime      int         `json:"datetime"`
		Cover         string      `json:"cover"`
		CoverWidth    int         `json:"cover_width"`
		CoverHeight   int         `json:"cover_height"`
		AccountURL    string      `json:"account_url"`
		AccountID     int         `json:"account_id"`
		Privacy       string      `json:"privacy"`
		Layout        string      `json:"layout"`
		Views         int         `json:"views"`
		Link          string      `json:"link"`
		Ups           int         `json:"ups"`
		Downs         int         `json:"downs"`
		Points        int         `json:"points"`
		Score         int         `json:"score"`
		IsAlbum       bool        `json:"is_album"`
		Vote          interface{} `json:"vote"`
		Favorite      bool        `json:"favorite"`
		Nsfw          bool        `json:"nsfw"`
		Section       string      `json:"section"`
		CommentCount  int         `json:"comment_count"`
		FavoriteCount int         `json:"favorite_count"`
		Topic         string      `json:"topic"`
		TopicID       int         `json:"topic_id"`
		ImagesCount   int         `json:"images_count"`
		InGallery     bool        `json:"in_gallery"`
		IsAd          bool        `json:"is_ad"`
		Tags          []struct {
			Name                   string      `json:"name"`
			DisplayName            string      `json:"display_name"`
			Followers              int         `json:"followers"`
			TotalItems             int         `json:"total_items"`
			Following              bool        `json:"following"`
			IsWhitelisted          bool        `json:"is_whitelisted"`
			BackgroundHash         string      `json:"background_hash"`
			ThumbnailHash          interface{} `json:"thumbnail_hash"`
			Accent                 string      `json:"accent"`
			BackgroundIsAnimated   bool        `json:"background_is_animated"`
			ThumbnailIsAnimated    bool        `json:"thumbnail_is_animated"`
			IsPromoted             bool        `json:"is_promoted"`
			Description            string      `json:"description"`
			LogoHash               interface{} `json:"logo_hash"`
			LogoDestinationURL     interface{} `json:"logo_destination_url"`
			DescriptionAnnotations struct {
			} `json:"description_annotations"`
		} `json:"tags"`
		AdType          int    `json:"ad_type"`
		AdURL           string `json:"ad_url"`
		InMostViral     bool   `json:"in_most_viral"`
		IncludeAlbumAds bool   `json:"include_album_ads"`
		Images          []struct {
			ID            string        `json:"id"`
			Title         interface{}   `json:"title"`
			Description   interface{}   `json:"description"`
			Datetime      int           `json:"datetime"`
			Type          string        `json:"type"`
			Animated      bool          `json:"animated"`
			Width         int           `json:"width"`
			Height        int           `json:"height"`
			Size          int           `json:"size"`
			Views         int           `json:"views"`
			Bandwidth     int64         `json:"bandwidth"`
			Vote          interface{}   `json:"vote"`
			Favorite      bool          `json:"favorite"`
			Nsfw          interface{}   `json:"nsfw"`
			Section       interface{}   `json:"section"`
			AccountURL    interface{}   `json:"account_url"`
			AccountID     interface{}   `json:"account_id"`
			IsAd          bool          `json:"is_ad"`
			InMostViral   bool          `json:"in_most_viral"`
			HasSound      bool          `json:"has_sound"`
			Tags          []interface{} `json:"tags"`
			AdType        int           `json:"ad_type"`
			AdURL         string        `json:"ad_url"`
			Edited        string        `json:"edited"`
			InGallery     bool          `json:"in_gallery"`
			Link          string        `json:"link"`
			CommentCount  interface{}   `json:"comment_count"`
			FavoriteCount interface{}   `json:"favorite_count"`
			Ups           interface{}   `json:"ups"`
			Downs         interface{}   `json:"downs"`
			Points        interface{}   `json:"points"`
			Score         interface{}   `json:"score"`
		} `json:"images"`
		AdConfig struct {
			SafeFlags       []string      `json:"safeFlags"`
			HighRiskFlags   []interface{} `json:"highRiskFlags"`
			UnsafeFlags     []interface{} `json:"unsafeFlags"`
			WallUnsafeFlags []interface{} `json:"wallUnsafeFlags"`
			ShowsAds        bool          `json:"showsAds"`
		} `json:"ad_config"`
	} `json:"data"`
	Success bool `json:"success"`
	Status  int  `json:"status"`
}

// Use the Imgur API to get a matching response struct for given category/resource
func getAPIResponse(category, id string) (ImgurResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.imgur.com/3/%s/%s", category, id), nil)
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}

	// Set headers
	imgurKey := os.Getenv("IMGUR_KEY")

	req.Header.Set("Authorization", fmt.Sprintf("Client-ID %s", imgurKey))

	client := &http.Client{Timeout: time.Second * 10}

	// Send request
	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	defer res.Body.Close()

	var apiResponse ImgurResponse
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&apiResponse)

	return apiResponse, err
}

// https://api.imgur.com/models/gallery_album
func imgurGallery(url, id string) (string, error) {
	apiResponse, err := getAPIResponse("gallery", id)
	if err != nil {
		return "", err
	}

	title := apiResponse.Data.Title
	if apiResponse.Data.ImagesCount > 1 {
		title = fmt.Sprintf("%s [%d images]", title, apiResponse.Data.ImagesCount)
	}
	if len(apiResponse.Data.Tags) > 0 {
		tags := []string{}
		for _, tag := range apiResponse.Data.Tags {
			tags = append(tags, tag.DisplayName)
		}
		title = fmt.Sprintf("%s [tags: %s]", title, strings.Join(tags, ", "))
	}

	return title, nil
}

// Just a normal album, not in the public gallery(?)
// https://api.imgur.com/models/album
func imgurAlbum(url, id string) (string, error) {
	apiResponse, err := getAPIResponse("album", id)
	if err != nil {
		return "", err
	}

	title := apiResponse.Data.Title
	if apiResponse.Data.ImagesCount > 1 {
		title = fmt.Sprintf("%s [%d images]", title, apiResponse.Data.ImagesCount)
	}
	if len(apiResponse.Data.Tags) > 0 {
		tags := []string{}
		for _, tag := range apiResponse.Data.Tags {
			tags = append(tags, tag.DisplayName)
		}
		title = fmt.Sprintf("%s [tags: %s]", title, strings.Join(tags, ", "))
	}

	return title, nil
}

// Imgur titles are always useless, just don't return anything
func Imgur(url string) (string, error) {
	match := galleryRegex.FindStringSubmatch(url)
	if len(match) > 0 {
		return imgurGallery(url, match[1])
	}

	match = albumRegex.FindStringSubmatch(url)
	if len(match) > 0 {
		return imgurAlbum(url, match[1])
	}

	// Nothing to be done
	return "", nil
}

// Register the handler function with corresponding regex
func init() {
	lambda.RegisterHandler(".*?imgur.com.*", Imgur)
}
