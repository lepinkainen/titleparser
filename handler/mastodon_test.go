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
		{"https://mastodon.social/@username", false},  // No post ID
		{"https://mastodon.social/users/username", false}, // Not /@username format
		{"https://twitter.com/username/status/123456789", false}, // Not Mastodon
		{"https://mastodon.social/about", false}, // Not a post
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