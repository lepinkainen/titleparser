//go:build !ci

package handler

import (
	"regexp"
	"testing"
)

func TestYoutube(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// Video tests
		{"Marvel 1", args{url: "https://www.youtube.com/watch?v=QdpxoFcdORI"}, `Marvel Studios Celebrates The Movies by Marvel Entertainment \[3m11s - \d+M views - \d+ (hours|days?|weeks?|months?|years?) ago\]`, false},
		{"Marvel 1 - short", args{url: "https://youtu.be/QdpxoFcdORI"}, `Marvel Studios Celebrates The Movies by Marvel Entertainment \[3m11s - \d+M views - \d+ (hours|days?|weeks?|months?|years?) ago\]`, false},
		{"Age restricted", args{url: "https://www.youtube.com/watch?v=EX_8ZjT2sO4"}, `Grenouer - Alone in the Dark - \[UNCENSORED - AGE RESTRICTED\] by GrenouerVEVO \[3m59s - \d+M views - \d years ago - age restricted\]`, false},
		{"At timestamp", args{url: "https://youtu.be/EX_8ZjT2sO4?t=98"}, `Grenouer - Alone in the Dark - \[UNCENSORED - AGE RESTRICTED\] by GrenouerVEVO \[3m59s - \d+M views - \d years ago - age restricted\]`, false},
		{"Gangnam style", args{url: "https://www.youtube.com/watch?v=9bZkp7q19f0"}, `PSY - GANGNAM STYLE\(강남스타일\) M/V by officialpsy \[4m13s - \d+Billion views - \d+ years ago\]`, false},

		// Channel tests
		{"Channel handle", args{url: "https://www.youtube.com/@GoogleDevelopers"}, `Google for Developers \[Channel - \d+[\w.]+ subscribers - \d+[\w.]+ videos - created \d+ years ago\]`, false},
		{"Channel custom", args{url: "https://www.youtube.com/c/GoogleDevelopers"}, `Google for Developers \[Channel - \d+[\w.]+ subscribers - \d+[\w.]+ videos - created \d+ years ago\]`, false},
		{"Channel ID", args{url: "https://www.youtube.com/channel/UC_x5XG1OV2P6uZZ5FSM9Ttw"}, `Google for Developers \[Channel - \d+[\w.]+ subscribers - \d+[\w.]+ videos - created \d+ years ago\]`, false},
		{"Channel user", args{url: "https://www.youtube.com/user/GoogleDevelopers"}, `Google for Developers \[Channel - \d+[\w.]+ subscribers - \d+[\w.]+ videos - created \d+ years ago\]`, false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := Youtube(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("Youtube() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			match, err := regexp.MatchString(tt.want, got)
			if err != nil || !match {
				t.Errorf("Youtube() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractVideoID(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{"Short URL", "https://youtu.be/QdpxoFcdORI", "QdpxoFcdORI"},
		{"Short URL with timestamp", "https://youtu.be/QdpxoFcdORI?t=98", "QdpxoFcdORI"},
		{"Long URL", "https://www.youtube.com/watch?v=QdpxoFcdORI", "QdpxoFcdORI"},
		{"Long URL with params", "https://www.youtube.com/watch?v=QdpxoFcdORI&t=98&list=abc", "QdpxoFcdORI"},
		{"Channel URL", "https://www.youtube.com/@GoogleDevelopers", ""},
		{"Invalid URL", "https://example.com", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractVideoID(tt.url); got != tt.want {
				t.Errorf("extractVideoID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractChannelInfo(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		wantID    string
		wantParam string
	}{
		{"Handle URL", "https://www.youtube.com/@GoogleDevelopers", "GoogleDevelopers", "forHandle"},
		{"Custom URL", "https://www.youtube.com/c/GoogleDevelopers", "GoogleDevelopers", "forUsername"},
		{"Channel ID URL", "https://www.youtube.com/channel/UC_x5XG1OV2P6uZZ5FSM9Ttw", "UC_x5XG1OV2P6uZZ5FSM9Ttw", "id"},
		{"User URL", "https://www.youtube.com/user/GoogleDevelopers", "GoogleDevelopers", "forUsername"},
		{"Direct URL", "https://www.youtube.com/GoogleDevelopers", "GoogleDevelopers", "forUsername"},
		{"Video URL", "https://www.youtube.com/watch?v=QdpxoFcdORI", "", ""},
		{"Invalid URL", "https://example.com", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotID, gotParam := ExtractChannelInfo(tt.url)
			if gotID != tt.wantID {
				t.Errorf("extractChannelInfo() gotID = %v, want %v", gotID, tt.wantID)
			}
			if gotParam != tt.wantParam {
				t.Errorf("extractChannelInfo() gotParam = %v, want %v", gotParam, tt.wantParam)
			}
		})
	}
}
