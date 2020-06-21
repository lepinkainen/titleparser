package handler

import "testing"

func TestHackerNews(t *testing.T) {
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
		{"Story 1", args{url: "https://news.ycombinator.com/item?id=23439437"}, "A List of Hacker News's Undocumented Features and Behaviors (2018) by billme [673 points]", false},
		{"Story 2", args{url: "https://news.ycombinator.com/item?id=23435805"}, "USB-C is still a mess by vo2maxer [216 points]", false},
		{"Wrong URL", args{url: "http://mantta.fi"}, "", false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := HackerNews(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("HackerNews() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("HackerNews() = %v, want %v", got, tt.want)
			}
		})
	}
}
