package handler

import (
	"regexp"
	"testing"
)

func TestReddit(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"Basic post", args{url: "https://www.reddit.com/r/funny/comments/np9b9b/for_those_having_trouble_finding_it/"}, `For those having trouble finding it. \[\d+ pts, \d+ comments, \d (days?|weeks?) ago]`, false},
		{"Basic post 2", args{url: "https://www.reddit.com/r/CryptoCurrency/comments/noztp7/binance_ceo_cz_shades_elon_musk_in_tweet_when_you/"}, `Binance CEO, CZ, shades Elon Musk in tweet. - ''When you use electricity to run cars, it’s environmentally friendly. When you use electricity to run the most efficient financial networks in the world, it’s an environmental concern.'' \[\d+ pts, \d+ comments, \d (days?|weeks?) ago]`, false},
		{"Gfycat post", args{url: "https://www.reddit.com/r/GifRecipes/comments/naqcu4/321_method_bbq_ribs/"}, `3-2-1 Method BBQ Ribs \[\d+ pts, \d+ comments, \d weeks ago]`, false},
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
			match, err := regexp.MatchString(tt.want, got)
			if err != nil || !match {
				t.Errorf("Reddit() = %v, want %v", got, tt.want)
			}
			match, err = regexp.MatchString(RedditMatch, tt.args.url)
			if err != nil || !match {
				t.Errorf("Reddit() URL didn't match regex")
			}
		})
	}
}
