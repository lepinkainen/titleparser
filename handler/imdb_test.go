//go:build !ci

package handler

import (
	"regexp"
	"testing"
)

// TODO: Make test use golden files instead of online testing
func TestOMDB(t *testing.T) {
	t.Parallel()

	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    *regexp.Regexp
		wantErr bool
	}{
		{"The Matrix - movie",
			args{url: "https://www.imdb.com/title/tt0133093/"},
			regexp.MustCompile(`^The Matrix \(1999\) \[IMDb \d\.\d/10\] \[RT \d+%\] \[Meta \d+/100\]$`),
			false},
		{"LOTR - movie - ref",
			args{url: "https://www.imdb.com/title/tt0120737/?ref_=fn_al_tt_1"},
			regexp.MustCompile(`^The Lord of the Rings: The Fellowship of the Ring \(2001\) \[IMDb \d\.\d/10\] \[RT \d+%\] \[Meta \d+/100\]$`),
			false},
		{"MacGyver - TV",
			args{url: "https://www.imdb.com/title/tt0088559/"},
			regexp.MustCompile(`^MacGyver \(1985.1992\) \[IMDb \d\.\d/10\] \[RT (N/A|\d+%)\] \[Meta (N/A|\d+/100)\]$`),
			false},
		{"No ID in URL",
			args{url: "https://www.imdb.com/"},
			nil,
			true},
		{"Megaforce - movie",
			args{url: "https://www.imdb.com/title/tt0084316/"},
			regexp.MustCompile(`^Megaforce \(1982\) \[IMDb \d\.\d/10\] \[RT \d+%\] \[Meta \d+/100\]$`),
			false},
		{"The Matrix - movie - no ending slash",
			args{url: "https://www.imdb.com/title/tt0133093"},
			regexp.MustCompile(`^The Matrix \(1999\) \[IMDb \d\.\d/10\] \[RT \d+%\] \[Meta \d+/100\]$`),
			false},
		{"Wrong URL", args{url: "http://mantta.fi"}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := OMDB(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("OMDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if !tt.want.MatchString(got) {
				t.Errorf("OMDB() = %v, want match %v", got, tt.want)
			}
		})
	}
}
