package handler

import "testing"

func TestVerkkokauppa(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"Nikko", args{url: "https://www.verkkokauppa.com/fi/product/40229/gxqht/Nikko-Vaporizr-2-kauko-ohjattava-auto-sininen"}, "Nikko Vaporizr 2 -kauko-ohjattava auto, sininen", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Verkkokauppa(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("Verkkokauppa() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Verkkokauppa() = %v, want %v", got, tt.want)
			}
		})
	}
}
