//go:build !ci

package handler

import (
	"regexp"
	"strings"
	"testing"
)

func TestTheRegister(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name        string
		args        args
		want        string
		wantErr     bool
		errContains string
	}{
		{"Basic article", args{url: "https://www.theregister.com/2022/03/21/google_messages_gdpr/"}, `Messages, Dialer apps sent text, call info to Google`, false, ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := TheRegister(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("TheRegister() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If we expect an error and have specific error text to check for
			if tt.wantErr && tt.errContains != "" && (err == nil || !strings.Contains(err.Error(), tt.errContains)) {
				t.Errorf("TheRegister() error = %v, should contain %v", err, tt.errContains)
				return
			}

			// If no error, check the content
			if err == nil {
				// For The Register, we just check that the title contains expected keywords
				// since titles can change over time
				if !strings.Contains(strings.ToLower(got), strings.ToLower(tt.want)) {
					t.Errorf("TheRegister() = '%v', want it to contain '%v'", got, tt.want)
				}
			}

			// Test that the URL matches our regex pattern
			match, err := regexp.MatchString(TheRegisterMatch, tt.args.url)
			if err != nil || !match {
				t.Errorf("TheRegister() URL '%s' didn't match regex '%s'", tt.args.url, TheRegisterMatch)
			}
		})
	}
}

func TestTheRegisterMatch(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want bool
	}{
		{"Basic theregister.com URL", "https://www.theregister.com/2022/03/21/google_messages_gdpr/", true},
		{"theregister.com without www", "https://theregister.com/2024/01/15/chrome_password_manager/", true},
		{"theregister.com subdomain", "https://feeds.theregister.com/rss/", true},
		{"Non-theregister URL", "https://www.reddit.com/r/programming/", false},
		{"Similar domain but not theregister", "https://www.nottheregister.com/article/", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := regexp.MatchString(TheRegisterMatch, tt.url)
			if err != nil {
				t.Errorf("regex error: %v", err)
				return
			}
			if match != tt.want {
				t.Errorf("TheRegisterMatch for '%s' = %v, want %v", tt.url, match, tt.want)
			}
		})
	}
}
