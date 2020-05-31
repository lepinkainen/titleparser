package handler

import "testing"

func TestYlilauta(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"Basic test", args{url: "https://ylilauta.org/arkisto/122796229"}, "Tässä langassa pukeudutaan kauppakasseihin - Arkisto | Ylilauta", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Ylilauta(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ylilauta() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Ylilauta() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetURL(t *testing.T) {
	type args struct {
		url     string
		authKey string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"Basic test", args{url: "https://ylilauta.org/arkisto/122796229", authKey: ""}, "Tässä langassa pukeudutaan kauppakasseihin - Arkisto | Ylilauta", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := handleYlilauta(tt.args.url, tt.args.authKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
