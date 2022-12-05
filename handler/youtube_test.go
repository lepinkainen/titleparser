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
		{"Marvel 1", args{url: "https://www.youtube.com/watch?v=QdpxoFcdORI"}, `Marvel Studios Celebrates The Movies by Marvel Entertainment \[3m11s - \d+M views - \d+ (hours|days?|weeks?|months?|years?) ago\]`, false},
		{"Marvel 1 - short", args{url: "https://youtu.be/QdpxoFcdORI"}, `Marvel Studios Celebrates The Movies by Marvel Entertainment \[3m11s - \d+M views - \d+ (hours|days?|weeks?|months?|years?) ago\]`, false},
		{"Age restricted", args{url: "https://www.youtube.com/watch?v=EX_8ZjT2sO4"}, `Grenouer - Alone in the Dark - \[UNCENSORED - AGE RESTRICTED\] by GrenouerVEVO \[3m59s - \d+M views - \d years ago\]`, false},
		{"At timestamp", args{url: "https://youtu.be/EX_8ZjT2sO4?t=98"}, `Grenouer - Alone in the Dark - \[UNCENSORED - AGE RESTRICTED\] by GrenouerVEVO \[3m59s - \d+M views - \d years ago\]`, false},
		{"Gangnam style", args{url: "https://www.youtube.com/watch?v=9bZkp7q19f0"}, `PSY - GANGNAM STYLE\(강남스타일\) M/V by officialpsy \[4m13s - \d+Billion views - \d+ years ago\]`, false},
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
