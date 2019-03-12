package handler

import "testing"

// TODO: Make test use golden files instead of online testing
func TestOMDB(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		/*
			{"The Matrix",
				args{url: "https://www.imdb.com/title/tt0133093/"},
				"The Matrix (1999) [IMDb 8.7/10] [RT 88%] [Meta 73/100]",
				false},
		*/
		{"No ID in URL",
			args{url: "https://www.imdb.com/"},
			"",
			true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := OMDB(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("OMDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("OMDB() = %v, want %v", got, tt.want)
			}
		})
	}
}
