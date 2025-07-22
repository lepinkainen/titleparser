package handler

import (
	"regexp"
	"testing"
)

func TestMastodonRegex(t *testing.T) {
	t.Parallel()

	// Test cases: URL and whether it should match
	testCases := []struct {
		url     string
		matches bool
	}{
		// Valid Mastodon URLs
		{"https://mastodon.social/@username/123456789", true},
		{"https://fosstodon.org/@someuser/987654321", true},
		{"https://hachyderm.io/@techuser/112233445566", true},
		{"https://infosec.exchange/@securityperson/111222333", true},

		// Invalid URLs
		{"https://mastodon.social/@username", false},             // No post ID
		{"https://mastodon.social/users/username", false},        // Not /@username format
		{"https://twitter.com/username/status/123456789", false}, // Not Mastodon
		{"https://mastodon.social/about", false},                 // Not a post
	}

	// Compile the regex once
	r := regexp.MustCompile(MastodonMatch)

	// Run the tests
	for _, tc := range testCases {
		match := r.MatchString(tc.url)
		if match != tc.matches {
			t.Errorf("URL %s: expected match=%v, got %v", tc.url, tc.matches, match)
		}
	}
}

func TestParseMastodonURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		url          string
		wantInstance string
		wantUsername string
		wantStatusID string
		wantErr      bool
	}{
		{
			name:         "Valid URL - mastodon.social",
			url:          "https://mastodon.social/@username/123456789",
			wantInstance: "mastodon.social",
			wantUsername: "username",
			wantStatusID: "123456789",
			wantErr:      false,
		},
		{
			name:         "Valid URL - other instance",
			url:          "https://fosstodon.org/@someuser/987654321",
			wantInstance: "fosstodon.org",
			wantUsername: "someuser",
			wantStatusID: "987654321",
			wantErr:      false,
		},
		{
			name:         "Invalid URL - no status ID",
			url:          "https://mastodon.social/@username",
			wantInstance: "",
			wantUsername: "",
			wantStatusID: "",
			wantErr:      true,
		},
		{
			name:         "Invalid URL - not Mastodon format",
			url:          "https://mastodon.social/users/username",
			wantInstance: "",
			wantUsername: "",
			wantStatusID: "",
			wantErr:      true,
		},
	}

	for _, tc := range testCases {
		tc := tc // Capture range variable for parallel execution
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			instance, username, statusID, err := parseMastodonURL(tc.url)

			if (err != nil) != tc.wantErr {
				t.Errorf("parseMastodonURL() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if !tc.wantErr {
				if instance != tc.wantInstance {
					t.Errorf("parseMastodonURL() instance = %v, want %v", instance, tc.wantInstance)
				}

				if username != tc.wantUsername {
					t.Errorf("parseMastodonURL() username = %v, want %v", username, tc.wantUsername)
				}

				if statusID != tc.wantStatusID {
					t.Errorf("parseMastodonURL() statusID = %v, want %v", statusID, tc.wantStatusID)
				}
			}
		})
	}
}

func TestStripHTML(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name:     "Basic HTML tags",
			html:     "<p>This is a <strong>test</strong> paragraph</p>",
			expected: "This is a test paragraph",
		},
		{
			name:     "With attributes",
			html:     "<div class=\"content\">Hello <span style=\"color:red\">world</span>!</div>",
			expected: "Hello world!",
		},
		{
			name:     "No HTML",
			html:     "Just plain text",
			expected: "Just plain text",
		},
		{
			name:     "Nested tags",
			html:     "<div><p>Nested <em><strong>content</strong></em></p></div>",
			expected: "Nested content",
		},
		{
			name:     "With extra spaces",
			html:     "  <p>  Trim   spaces  </p>  ",
			expected: "Trim   spaces",
		},
	}

	for _, tc := range testCases {
		tc := tc // Capture range variable for parallel execution
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result := stripHTML(tc.html)

			if result != tc.expected {
				t.Errorf("stripHTML() = %q, want %q", result, tc.expected)
			}
		})
	}
}

func TestExtractOpenGraphTitle(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name:     "Valid OG title",
			html:     "<html><head><meta property=\"og:title\" content=\"This is the OG title\"></head><body></body></html>",
			expected: "This is the OG title",
		},
		{
			name:     "No OG title",
			html:     "<html><head><title>Page Title</title></head><body></body></html>",
			expected: "",
		},
		{
			name:     "Empty OG title",
			html:     "<html><head><meta property=\"og:title\" content=\"\"></head><body></body></html>",
			expected: "",
		},
	}

	for _, tc := range testCases {
		tc := tc // Capture range variable for parallel execution
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result := extractOpenGraphTitle(tc.html)

			if result != tc.expected {
				t.Errorf("extractOpenGraphTitle() = %q, want %q", result, tc.expected)
			}
		})
	}
}

func TestExtractTitle(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name:     "Valid title",
			html:     "<html><head><title>Page Title</title></head><body></body></html>",
			expected: "Page Title",
		},
		{
			name:     "No title",
			html:     "<html><head></head><body></body></html>",
			expected: "",
		},
		{
			name:     "Empty title",
			html:     "<html><head><title></title></head><body></body></html>",
			expected: "",
		},
	}

	for _, tc := range testCases {
		tc := tc // Capture range variable for parallel execution
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result := extractTitle(tc.html)

			if result != tc.expected {
				t.Errorf("extractTitle() = %q, want %q", result, tc.expected)
			}
		})
	}
}
