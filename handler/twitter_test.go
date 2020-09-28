package handler

import "testing"
import "strings"

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
		{"Kohli 1 - emoji,link", args{url: "https://twitter.com/RahulKohli13/status/1263946250077929473"}, "Rahul Kohli (✔@RahulKohli13) 2m: Becoming an actor can be incredibly challenging and extremely daunting but if you really want get off to a great start, stop using the fucking hashtag #ActorsLife on Instagram with your stupid fucking pictures of you posing at Starbucks, “prepping for your audition 💅🏽”. https://t.co/yeTkJkBx2Z", false},
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
			if !strings.starsWith(got, tt.want) {
				t.Errorf("Twitter() = %v, want %v", got, tt.want)
			}
		})
	}
}
