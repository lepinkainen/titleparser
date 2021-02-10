package handler

import "testing"

func TestImgur(t *testing.T) {
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
		{"gallery test", args{url: "https://imgur.com/gallery/md2Sxjm"}, "Chill in the face of aggression. [tags: protest, take that, current events, police brutality, america]", false},
		{"gallery test 2", args{url: "https://imgur.com/gallery/QXR4OL2"}, "Another Soulsborne dump? Why not! [44 images] [tags: darksouls, bloodborne, Dark Souls, dark souls 3, dark souls 2 is actually ok]", false},
		{"gallery test 3", args{url: "https://imgur.com/gallery/YCNzrKx"}, "As pertinent and poignant as ever. [tags: reaction, current events]", false},
		{"album test", args{url: "https://imgur.com/a/X2PcObK"}, "Assembly 1998 prizegivign. [2 images]", false},
		{"album test/no title", args{url: "https://imgur.com/a/MZY7mkE"}, "", false},
		{"image test", args{url: "https://imgur.com/BGMckfX"}, "Ella Love", false},
		{"subreddit image test", args{url: "https://i.imgur.com/k3w8kHG.jpg"}, "Arkiliikenne [/r/Suomi]", false},
		{"subreddit image test direct 1", args{url: "https://imgur.com/k3w8kHG"}, "Arkiliikenne [/r/Suomi]", false},
		{"subreddit image test 2", args{url: "https://i.imgur.com/dJJbwhM.jpg"}, "Stay safe out there. [/r/MTB]", false},
		{"subreddit image test direct 2", args{url: "https://imgur.com/dJJbwhM"}, "Stay safe out there. [/r/MTB]", false},
		{"subreddit gifv test", args{url: "https://i.imgur.com/OiocRjL.gifv"}, "The correct usage of a phone [/r/gifs]", false},
		{"tag gallery image", args{url: "https://imgur.com/t/funny/MWvY6dD"}, "Isolation day:....lost count [tags: armpit, mildly interesting, the more you know, funny, fart]", false},
		{"tag gallery image 2", args{url: "https://imgur.com/t/fart/srK588K"}, "Sleepy Rue [tags: doggo, daggs, happy friday, sweet jesus, fart]", false},
		{"Wrong URL", args{url: "http://mantta.fi"}, "", false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
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
