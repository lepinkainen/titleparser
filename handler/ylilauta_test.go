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
		{"Archive URL test", args{url: "https://ylilauta.org/arkisto/122796229"}, "Tässä langassa pukeudutaan kauppakasseihin - Arkisto", false},
		{"Active URL test", args{url: "https://ylilauta.org/sekalainen/125520689"}, "Vassari yritti ryöstää rekan mutta jäikin rekan alle:D - Suvaitsevainen", false},
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
