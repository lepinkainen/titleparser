package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/lepinkainen/titleparser/lambda"
	log "github.com/sirupsen/logrus"
)

var RedditMatch = ".*reddit.com/r/.*/comments/.*/.*"

type RedditPost []struct {
	Kind string `json:"kind"`
	Data struct {
		Modhash  string `json:"modhash"`
		Dist     int    `json:"dist"`
		Children []struct {
			Kind string `json:"kind"`
			Data struct {
				ApprovedAtUtc              any     `json:"approved_at_utc"`
				Subreddit                  string  `json:"subreddit"`
				Selftext                   string  `json:"selftext"`
				UserReports                []any   `json:"user_reports"`
				Saved                      bool    `json:"saved"`
				ModReasonTitle             any     `json:"mod_reason_title"`
				Gilded                     int     `json:"gilded"`
				Clicked                    bool    `json:"clicked"`
				Title                      string  `json:"title"`
				LinkFlairRichtext          []any   `json:"link_flair_richtext"`
				SubredditNamePrefixed      string  `json:"subreddit_name_prefixed"`
				Hidden                     bool    `json:"hidden"`
				Pwls                       int     `json:"pwls"`
				LinkFlairCSSClass          any     `json:"link_flair_css_class"`
				Downs                      int     `json:"downs"`
				TopAwardedType             any     `json:"top_awarded_type"`
				ParentWhitelistStatus      string  `json:"parent_whitelist_status"`
				HideScore                  bool    `json:"hide_score"`
				Name                       string  `json:"name"`
				Quarantine                 bool    `json:"quarantine"`
				LinkFlairTextColor         string  `json:"link_flair_text_color"`
				UpvoteRatio                float64 `json:"upvote_ratio"`
				AuthorFlairBackgroundColor any     `json:"author_flair_background_color"`
				SubredditType              string  `json:"subreddit_type"`
				Ups                        int     `json:"ups"`
				TotalAwardsReceived        int     `json:"total_awards_received"`
				MediaEmbed                 struct {
				} `json:"media_embed"`
				AuthorFlairTemplateID any    `json:"author_flair_template_id"`
				IsOriginalContent     bool   `json:"is_original_content"`
				AuthorFullname        string `json:"author_fullname"`
				SecureMedia           any    `json:"secure_media"`
				IsRedditMediaDomain   bool   `json:"is_reddit_media_domain"`
				IsMeta                bool   `json:"is_meta"`
				Category              any    `json:"category"`
				SecureMediaEmbed      struct {
				} `json:"secure_media_embed"`
				LinkFlairText       any    `json:"link_flair_text"`
				CanModPost          bool   `json:"can_mod_post"`
				Score               int    `json:"score"`
				ApprovedBy          any    `json:"approved_by"`
				AuthorPremium       bool   `json:"author_premium"`
				Thumbnail           string `json:"thumbnail"`
				Edited              any    `json:"edited"` // boolean or float timestamp
				AuthorFlairCSSClass any    `json:"author_flair_css_class"`
				AuthorFlairRichtext []any  `json:"author_flair_richtext"`
				Gildings            struct {
					Gid1 int `json:"gid_1"`
				} `json:"gildings"`
				PostHint            string  `json:"post_hint"`
				ContentCategories   any     `json:"content_categories"`
				IsSelf              bool    `json:"is_self"`
				ModNote             any     `json:"mod_note"`
				Created             float64 `json:"created"`
				LinkFlairType       string  `json:"link_flair_type"`
				Wls                 int     `json:"wls"`
				RemovedByCategory   any     `json:"removed_by_category"`
				BannedBy            any     `json:"banned_by"`
				AuthorFlairType     string  `json:"author_flair_type"`
				Domain              string  `json:"domain"`
				AllowLiveComments   bool    `json:"allow_live_comments"`
				SelftextHTML        any     `json:"selftext_html"`
				Likes               any     `json:"likes"`
				SuggestedSort       any     `json:"suggested_sort"`
				BannedAtUtc         any     `json:"banned_at_utc"`
				URLOverriddenByDest string  `json:"url_overridden_by_dest"`
				ViewCount           any     `json:"view_count"`
				Archived            bool    `json:"archived"`
				NoFollow            bool    `json:"no_follow"`
				IsCrosspostable     bool    `json:"is_crosspostable"`
				Pinned              bool    `json:"pinned"`
				Over18              bool    `json:"over_18"`
				Preview             struct {
					Images []struct {
						Source struct {
							URL    string `json:"url"`
							Width  int    `json:"width"`
							Height int    `json:"height"`
						} `json:"source"`
						Resolutions []struct {
							URL    string `json:"url"`
							Width  int    `json:"width"`
							Height int    `json:"height"`
						} `json:"resolutions"`
						Variants struct {
						} `json:"variants"`
						ID string `json:"id"`
					} `json:"images"`
					Enabled bool `json:"enabled"`
				} `json:"preview"`
				AllAwardings []struct {
					GiverCoinReward          any    `json:"giver_coin_reward"`
					SubredditID              any    `json:"subreddit_id"`
					IsNew                    bool   `json:"is_new"`
					DaysOfDripExtension      int    `json:"days_of_drip_extension"`
					CoinPrice                int    `json:"coin_price"`
					ID                       string `json:"id"`
					PennyDonate              any    `json:"penny_donate"`
					CoinReward               int    `json:"coin_reward"`
					IconURL                  string `json:"icon_url"`
					DaysOfPremium            int    `json:"days_of_premium"`
					IconHeight               int    `json:"icon_height"`
					TiersByRequiredAwardings any    `json:"tiers_by_required_awardings"`
					ResizedIcons             []struct {
						URL    string `json:"url"`
						Width  int    `json:"width"`
						Height int    `json:"height"`
					} `json:"resized_icons"`
					IconWidth                        int    `json:"icon_width"`
					StaticIconWidth                  int    `json:"static_icon_width"`
					StartDate                        any    `json:"start_date"`
					IsEnabled                        bool   `json:"is_enabled"`
					AwardingsRequiredToGrantBenefits any    `json:"awardings_required_to_grant_benefits"`
					Description                      string `json:"description"`
					EndDate                          any    `json:"end_date"`
					SubredditCoinReward              int    `json:"subreddit_coin_reward"`
					Count                            int    `json:"count"`
					StaticIconHeight                 int    `json:"static_icon_height"`
					Name                             string `json:"name"`
					ResizedStaticIcons               []struct {
						URL    string `json:"url"`
						Width  int    `json:"width"`
						Height int    `json:"height"`
					} `json:"resized_static_icons"`
					IconFormat    any    `json:"icon_format"`
					AwardSubType  string `json:"award_sub_type"`
					PennyPrice    any    `json:"penny_price"`
					AwardType     string `json:"award_type"`
					StaticIconURL string `json:"static_icon_url"`
				} `json:"all_awardings"`
				Awarders                 []any   `json:"awarders"`
				MediaOnly                bool    `json:"media_only"`
				CanGild                  bool    `json:"can_gild"`
				Spoiler                  bool    `json:"spoiler"`
				Locked                   bool    `json:"locked"`
				AuthorFlairText          any     `json:"author_flair_text"`
				TreatmentTags            []any   `json:"treatment_tags"`
				Visited                  bool    `json:"visited"`
				RemovedBy                any     `json:"removed_by"`
				NumReports               any     `json:"num_reports"`
				Distinguished            any     `json:"distinguished"`
				SubredditID              string  `json:"subreddit_id"`
				ModReasonBy              any     `json:"mod_reason_by"`
				RemovalReason            any     `json:"removal_reason"`
				LinkFlairBackgroundColor string  `json:"link_flair_background_color"`
				ID                       string  `json:"id"`
				IsRobotIndexable         bool    `json:"is_robot_indexable"`
				NumDuplicates            int     `json:"num_duplicates"`
				ReportReasons            any     `json:"report_reasons"`
				Author                   string  `json:"author"`
				DiscussionType           any     `json:"discussion_type"`
				NumComments              int     `json:"num_comments"`
				SendReplies              bool    `json:"send_replies"`
				Media                    any     `json:"media"`
				Permalink                string  `json:"permalink"`
				WhitelistStatus          string  `json:"whitelist_status"`
				Stickied                 bool    `json:"stickied"`
				URL                      string  `json:"url"`
				SubredditSubscribers     int     `json:"subreddit_subscribers"`
				CreatedUtc               float64 `json:"created_utc"`
				NumCrossposts            int     `json:"num_crossposts"`
				ModReports               []any   `json:"mod_reports"`
				IsVideo                  bool    `json:"is_video"`
			} `json:"data"`
		} `json:"children"`
		After  any `json:"after"`
		Before any `json:"before"`
	} `json:"data"`
}

func Reddit(url string) (string, error) {
	if strings.HasSuffix(url, "/") {
		url = fmt.Sprintf("%s.json", url)
	} else {
		url = fmt.Sprintf("%s/.json", url)
	}

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

	var apiResponse RedditPost
	dec := json.NewDecoder(res.Body)
	//bytes, _ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(bytes))
	err = dec.Decode(&apiResponse)
	if err != nil {
		_ = fmt.Errorf("error decoding API response: %v", err)
		return "", err
	}

	//fmt.Println(apiResponse[0].Data.Children[0])

	data := apiResponse[0].Data.Children[0].Data
	over_18 := data.Over18
	created := time.Unix(int64(data.CreatedUtc), 0)
	age := humanize.RelTime(created, time.Now(), "ago", "")
	//author := data.Author

	title := fmt.Sprintf("%s [%d pts, %d comments, %s]", data.Title, data.Score, data.NumComments, age)
	if over_18 {
		title = fmt.Sprintf("%s (NSFW)", title)
	}

	return title, nil
}

func init() {
	lambda.RegisterHandler(RedditMatch, Reddit)
}
