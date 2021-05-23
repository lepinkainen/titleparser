package handler

import "testing"

// TODO: Make test use golden files instead of online testing
func TestOMDB(t *testing.T) {
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
		{"The Matrix - movie",
			args{url: "https://www.imdb.com/title/tt0133093/"},
			"The Matrix (1999) [IMDb 8.7/10] [RT 88%] [Meta 73/100]",
			false},
		{"LOTR - movie - ref",
			args{url: "https://www.imdb.com/title/tt0120737/?ref_=fn_al_tt_1"},
			"The Lord of the Rings: The Fellowship of the Ring (2001) [IMDb 8.8/10] [RT 91%] [Meta 92/100]",
			false},
		{"MacGyver - TV",
			args{url: "https://www.imdb.com/title/tt0088559/"},
			"MacGyver (1985–1992) [IMDb 7.7/10] [RT N/A] [Meta N/A]",
			false},
		{"No ID in URL",
			args{url: "https://www.imdb.com/"},
			"",
			true},
		{"The Matrix - movie - no ending slash",
			args{url: "https://www.imdb.com/title/tt0133093"},
			"The Matrix (1999) [IMDb 8.7/10] [RT 88%] [Meta 73/100]",
			false},
		{"Wrong URL", args{url: "http://mantta.fi"}, "", true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
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
