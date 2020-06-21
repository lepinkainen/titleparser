package handler

import "testing"

func TestImgur(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"gallery test", args{url: "https://imgur.com/gallery/md2Sxjm"}, "Chill in the face of aggression. [tags: america, police brutality, protest, current events, take that]", false},
		{"gallery test 2", args{url: "https://imgur.com/gallery/QXR4OL2"}, "Another Soulsborne dump? Why not! [44 images] [tags: Dark Souls, bloodborne, darksouls, dark souls 3, dark souls 2 is actually ok]", false},
		{"gallery test 3", args{url: "https://imgur.com/gallery/YCNzrKx"}, "As pertinent and poignant as ever. [tags: reaction, current events]", false},
		{"album test", args{url: "https://imgur.com/a/X2PcObK"}, "Assembly 1998 prizegivign. [2 images]", false},
		{"album test/no title", args{url: "https://imgur.com/a/MZY7mkE"}, "", false},
		{"image test", args{url: "https://imgur.com/BGMckfX"}, "Ella Love", false},
		{"subreddit image test", args{url: "https://i.imgur.com/k3w8kHG.jpg"}, "Arkiliikenne [/r/Suomi]", false},
		{"subreddit image test direct 1", args{url: "https://imgur.com/k3w8kHG"}, "Arkiliikenne [/r/Suomi]", false},
		{"subreddit image test 2", args{url: "https://i.imgur.com/dJJbwhM.jpg"}, "Stay safe out there. [/r/MTB]", false},
		{"subreddit image test direct 2", args{url: "https://imgur.com/dJJbwhM"}, "Stay safe out there. [/r/MTB]", false},
		{"subreddit gifv test", args{url: "https://i.imgur.com/OiocRjL.gifv"}, "The correct usage of a phone [/r/gifs]", false},
		{"tag gallery image", args{url: "https://imgur.com/t/funny/MWvY6dD"}, "Isolation day:....lost count", false},
		{"Wrong URL", args{url: "http://mantta.fi"}, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Imgur(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("Imgur() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Imgur() = %v, want %v", got, tt.want)
			}
		})
	}
}
