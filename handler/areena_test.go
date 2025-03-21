package handler

import (
	"regexp"
	"testing"
)

func TestYleAreena(t *testing.T) {
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
		{"Old movie", args{url: "https://areena.yle.fi/1-4192173"}, "Florence", false},
		{"Old podcast", args{url: "https://areena.yle.fi/audio/1-1792200"}, "Perttu Häkkinen", false},
		{"Series main page", args{url: "https://areena.yle.fi/1-3371178"}, "Pikku Kakkonen", false},
		{"Series episode", args{url: "https://areena.yle.fi/1-50696546"}, "Maanantai 5.4.2021 | Pikku Kakkonen", false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := YleAreena(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("YleAreena() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			match, err := regexp.MatchString(tt.want, got)
			if err != nil || !match {
				t.Errorf("YleAreena() = %v, want %v", got, tt.want)
			}
		})
	}
}
