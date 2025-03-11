package handler

import (
	"regexp"
	"strings"
	"testing"
)

func TestReddit(t *testing.T) {
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
		{"Basic post", args{url: "https://www.reddit.com/r/funny/comments/np9b9b/for_those_having_trouble_finding_it/"}, `For those having trouble finding it.`, false, ""},
		{"Basic post 2", args{url: "https://www.reddit.com/r/selfhosted/comments/151clya/selfhosted_alternatives_for_dropbox/"}, `Selfhosted alternatives for Dropbox`, false, ""},
		{"Basic post 3", args{url: "https://www.reddit.com/r/MandelaEffect/comments/1ah7asy/i_think_i_figured_out_the_criminal_emoji_mandela/"}, `I think I figured out the criminal emoji Mandela effect.`, false, ""},
		{"Gfycat post", args{url: "https://www.reddit.com/r/GifRecipes/comments/naqcu4/321_method_bbq_ribs/"}, `3-2-1 Method BBQ Ribs`, false, ""},
		{"Nextfuckinglevel post", args{url: "https://www.reddit.com/r/nextfuckinglevel/comments/1j8vcog/finnish_freediver_olavi_paananen_broke_the_world/"}, `Finnish freediver olavi paananen, broke the world record diving 107 meters under the ice without flippers and wearing only swimming shorts`, false, ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := Reddit(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("Reddit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If we expect an error and have specific error text to check for
			if tt.wantErr && tt.errContains != "" && (err == nil || !strings.Contains(err.Error(), tt.errContains)) {
				t.Errorf("Reddit() error = %v, should contain %v", err, tt.errContains)
				return
			}

			tt.want = tt.want + ` \[\d+ pts, \d+ comments, \d (hours?|days?|weeks?|months?|years?) ago]`

			// If no error, check the content
			if err == nil {
				match, err := regexp.MatchString(tt.want, got)
				if err != nil || !match {
					t.Errorf("Reddit() = '%v', want '%v'", got, tt.want)
				}
			}

			match, err := regexp.MatchString(RedditMatch, tt.args.url)
			if err != nil || !match {
				t.Errorf("Reddit() URL didn't match regex")
			}
		})
	}
}
