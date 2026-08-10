//go:build !ci

package handler

import (
	"regexp"
	"strings"
	"testing"
)

func TestThreads(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"Share link", args{url: "https://www.threads.com/share/_0ozSh-x9/"}, "on Threads", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := Threads(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("Threads() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				// Title is "<author> (@handle) on Threads"; assert the stable
				// suffix rather than the author name which can change.
				if !strings.Contains(got, tt.want) {
					t.Errorf("Threads() = '%v', want it to contain '%v'", got, tt.want)
				}
				if strings.EqualFold(strings.TrimSpace(got), "Threads") {
					t.Errorf("Threads() = '%v', got the bare shell title (crawler UA not honoured?)", got)
				}
			}

			match, err := regexp.MatchString(ThreadsMatch, tt.args.url)
			if err != nil || !match {
				t.Errorf("Threads() URL '%s' didn't match regex '%s'", tt.args.url, ThreadsMatch)
			}
		})
	}
}

func TestThreadsMatch(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want bool
	}{
		{"Share link", "https://www.threads.com/share/_0ozSh-x9/", true},
		{"Post URL", "https://www.threads.com/@grimmemento/post/Db2xD9kiAFI", true},
		{"Legacy threads.net post", "https://www.threads.net/@someone/post/ABC123", true},
		{"Short /t/ link", "https://www.threads.com/t/ABC123/", true},
		{"Threads homepage", "https://www.threads.com/", false},
		{"Non-threads URL", "https://www.reddit.com/r/programming/", false},
		{"Similar domain but not threads", "https://www.notthreads.com/@x/post/1", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := regexp.MatchString(ThreadsMatch, tt.url)
			if err != nil {
				t.Errorf("regex error: %v", err)
				return
			}
			if match != tt.want {
				t.Errorf("ThreadsMatch for '%s' = %v, want %v", tt.url, match, tt.want)
			}
		})
	}
}
