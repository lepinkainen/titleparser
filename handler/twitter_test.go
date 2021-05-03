package handler

import (
	"fmt"
	"regexp"
	"testing"
)

func TestTwitter(t *testing.T) {
	t.Parallel()

	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"Kohli 1 - emoji,link", args{url: "https://twitter.com/RahulKohli13/status/1263946250077929473"}, `Rahul Kohli \(✔@RahulKohli13\) \d+(d|m|y): Becoming an actor can be incredibly challenging and extremely daunting but if you really want get off to a great start, stop using the fucking hashtag #ActorsLife on Instagram with your stupid fucking pictures of you posing at Starbucks, “prepping for your audition 💅🏽”. https://t.co/yeTkJkBx2Z \[. \d+ . \d+\]`, false},
		{"Foone 1 - image", args{url: "https://twitter.com/Foone/status/1389076911138115587"}, `foone \(@Foone\) \d+(m|h|d|y): and then I flashed it back. This thing is magic. https://t.co/4n8UxvEhvG \[. \d+ . \d+\]`, false},
		{"Foone 2 - video", args{url: "https://twitter.com/Foone/status/1250937428518436864"}, `foone \(@Foone\) \d+(m|h|d|y): So the show Shuriken Sentai Ninninger, the 39th show in the Super Sentai series (which formed the basis for the Power Rangers shows), had a monster named Yokai Mokumokuren, who is a supernaturally-infected computer keyboard. https://t.co/224u1P9iqb`, false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := Twitter(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("Twitter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			match, err := regexp.MatchString(tt.want, got)
			fmt.Println(match)
			if err != nil || !match {
				t.Errorf("Twitter() = %v, want %v", got, tt.want)
			}
		})
	}
}
