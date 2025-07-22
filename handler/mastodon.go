package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/lepinkainen/titleparser/lambda"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// MastodonMatch matches URLs with the Mastodon post pattern: /@username/numeric_id
// This pattern works across all Mastodon instances
var MastodonMatch = ".*/@[^/]+/[0-9]+"

// MastodonStatus represents a Mastodon status (post)
type MastodonStatus struct {
	ID               string                    `json:"id"`
	CreatedAt        time.Time                 `json:"created_at"`
	Content          string                    `json:"content"`
	ReblogsCount     int                       `json:"reblogs_count"`
	FavoritesCount   int                       `json:"favourites_count"`
	RepliesCount     int                       `json:"replies_count"`
	URL              string                    `json:"url"`
	Visibility       string                    `json:"visibility"`
	Language         string                    `json:"language"`
	Sensitive        bool                      `json:"sensitive"`
	Spoiler          string                    `json:"spoiler_text"`
	MediaAttachments []MastodonMediaAttachment `json:"media_attachments"`
	Account          MastodonAccount           `json:"account"`
	Application      MastodonApplication       `json:"application"`
	Mentions         []MastodonMention         `json:"mentions"`
	Tags             []MastodonTag             `json:"tags"`
	Card             *MastodonCard             `json:"card"`
}

// MastodonAccount represents a Mastodon user account
type MastodonAccount struct {
	ID             string `json:"id"`
	Username       string `json:"username"`
	Acct           string `json:"acct"`
	DisplayName    string `json:"display_name"`
	FollowersCount int    `json:"followers_count"`
	FollowingCount int    `json:"following_count"`
	StatusesCount  int    `json:"statuses_count"`
	Note           string `json:"note"`
	URL            string `json:"url"`
	Avatar         string `json:"avatar"`
}

// MastodonMediaAttachment represents media attached to a post
type MastodonMediaAttachment struct {
	ID          string `json:"id"`
	Type        string `json:"type"` // image, video, gifv, audio
	URL         string `json:"url"`
	PreviewURL  string `json:"preview_url"`
	Description string `json:"description"`
	Blurhash    string `json:"blurhash"`
}

// MastodonApplication represents the application that created the status
type MastodonApplication struct {
	Name    string `json:"name"`
	Website string `json:"website"`
}

// MastodonMention represents a mention in a post
type MastodonMention struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Acct     string `json:"acct"`
	URL      string `json:"url"`
}

// MastodonTag represents a hashtag in a post
type MastodonTag struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// MastodonCard represents a link preview card
type MastodonCard struct {
	URL          string `json:"url"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Type         string `json:"type"`
	Image        string `json:"image"`
	AuthorName   string `json:"author_name"`
	AuthorURL    string `json:"author_url"`
	ProviderName string `json:"provider_name"`
	ProviderURL  string `json:"provider_url"`
}

// parseMastodonURL extracts the instance domain, username, and status ID from a Mastodon URL
func parseMastodonURL(postURL string) (string, string, string, error) {
	// Parse the URL
	u, err := url.Parse(postURL)
	if err != nil {
		return "", "", "", errors.Wrap(err, "failed to parse URL")
	}

	// Extract the instance domain
	instance := u.Host

	// Extract username and status ID using regex
	re := regexp.MustCompile(`.*/@([^/]+)/([0-9]+)`)
	matches := re.FindStringSubmatch(postURL)
	if len(matches) < 3 {
		return "", "", "", errors.New("invalid Mastodon status URL format")
	}

	username := matches[1]
	statusID := matches[2]

	return instance, username, statusID, nil
}

// fetchStatusInfo fetches status information from the Mastodon API
func fetchStatusInfo(instance, statusID string) (*MastodonStatus, error) {
	// Construct the API URL
	apiURL := fmt.Sprintf("https://%s/api/v1/statuses/%s", instance, statusID)

	// Create a request with proper headers
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating request")
	}

	// Set headers
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Accept", "application/json")

	// Set client timeout
	client := &http.Client{Timeout: time.Second * 10}

	// Send request
	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "error sending request")
	}
	defer res.Body.Close()

	// Check if the request was successful
	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("API returned non-OK status: %d", res.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error reading response body")
	}

	// Parse the JSON
	var status MastodonStatus
	if err := json.Unmarshal(body, &status); err != nil {
		return nil, errors.Wrap(err, "error parsing JSON")
	}

	return &status, nil
}

// stripHTML removes HTML tags from a string
func stripHTML(html string) string {
	// Simple but effective for our needs
	re := regexp.MustCompile("<[^>]*>")
	text := re.ReplaceAllString(html, "")
	return strings.TrimSpace(text)
}

// Mastodon extracts information from a Mastodon post URL using the API
func Mastodon(url string) (string, error) {
	// Parse the Mastodon URL
	instance, _, statusID, err := parseMastodonURL(url)
	if err != nil {
		log.WithError(err).Error("Failed to parse Mastodon URL")
		return fallbackToScraping(url)
	}

	// Fetch status information
	status, err := fetchStatusInfo(instance, statusID)
	if err != nil {
		log.WithError(err).Error("Failed to fetch status info from API")
		return fallbackToScraping(url)
	}

	// Process the content
	content := stripHTML(status.Content)

	// Truncate content if too long
	if len(content) > 100 {
		content = content[:97] + "..."
	}

	// Format the relative time
	timestamp := humanize.RelTime(status.CreatedAt, time.Now(), "ago", "")

	// Build title with rich information
	title := fmt.Sprintf("@%s", status.Account.DisplayName)
	if status.Account.DisplayName != status.Account.Username {
		title += fmt.Sprintf(" (@%s)", status.Account.Username)
	}

	title += fmt.Sprintf(": %s", content)

	// Add engagement info
	engagement := []string{}
	if status.ReblogsCount > 0 {
		engagement = append(engagement, fmt.Sprintf("%d boosts", status.ReblogsCount))
	}
	if status.FavoritesCount > 0 {
		engagement = append(engagement, fmt.Sprintf("%d favs", status.FavoritesCount))
	}
	if status.RepliesCount > 0 {
		engagement = append(engagement, fmt.Sprintf("%d replies", status.RepliesCount))
	}

	// Add media info
	mediaInfo := ""
	if len(status.MediaAttachments) > 0 {
		types := make(map[string]int)
		for _, media := range status.MediaAttachments {
			types[media.Type]++
		}

		mediaLabels := []string{}
		for mediaType, count := range types {
			if count > 1 {
				mediaLabels = append(mediaLabels, fmt.Sprintf("%d %ss", count, mediaType))
			} else {
				mediaLabels = append(mediaLabels, mediaType)
			}
		}

		if len(mediaLabels) > 0 {
			mediaInfo = strings.Join(mediaLabels, ", ")
		}
	}

	// Format the final title
	if len(engagement) > 0 {
		title += fmt.Sprintf(" [%s]", strings.Join(engagement, ", "))
	}

	if mediaInfo != "" {
		title += fmt.Sprintf(" [Media: %s]", mediaInfo)
	}

	title += fmt.Sprintf(" [%s]", timestamp)

	// Add language if available and not default
	if status.Language != "" && status.Language != "en" {
		title += fmt.Sprintf(" [%s]", strings.ToUpper(status.Language))
	}

	return title, nil
}

// fallbackToScraping falls back to the old HTML scraping method if the API call fails
func fallbackToScraping(url string) (string, error) {
	// Create a request with proper headers
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error("Error creating request: ", err)
		return "", err
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
		log.Error("Error sending request: ", err)
		return "", err
	}
	defer res.Body.Close()

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error("Error reading response body: ", err)
		return "", err
	}

	// Try to find the OpenGraph title
	ogTitle := extractOpenGraphTitle(string(body))
	if ogTitle != "" {
		return ogTitle, nil
	}

	// Extract just the page title as a last resort
	title := extractTitle(string(body))
	if title != "" {
		return title, nil
	}

	return "Mastodon Post", nil
}

// extractOpenGraphTitle extracts the OpenGraph title from HTML
func extractOpenGraphTitle(html string) string {
	re := regexp.MustCompile(`<meta\s+property="og:title"\s+content="([^"]+)"`)
	matches := re.FindStringSubmatch(html)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// extractTitle extracts the page title from HTML
func extractTitle(html string) string {
	re := regexp.MustCompile(`<title>([^<]+)</title>`)
	matches := re.FindStringSubmatch(html)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func init() {
	lambda.RegisterHandler(MastodonMatch, Mastodon)
}
