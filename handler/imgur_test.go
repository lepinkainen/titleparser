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
