package handler

import "testing"

func TestYleAreena(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"Custom handler", args{url: "https://areena.yle.fi/1-4192173"}, "Areena custom handler", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := YleAreena(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("YleAreena() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("YleAreena() = %v, want %v", got, tt.want)
			}
		})
	}
}
