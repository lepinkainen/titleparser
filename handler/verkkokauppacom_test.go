package handler

import "testing"

func TestVerkkokauppa(t *testing.T) {
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
		{"Nikko", args{url: "https://www.verkkokauppa.com/fi/product/40229/gxqht/Nikko-Vaporizr-2-kauko-ohjattava-auto-sininen"}, "Nikko Vaporizr 2 -kauko-ohjattava auto, sininen", false},
		{"Fujtech", args{url: "https://www.verkkokauppa.com/fi/product/46243/msgsg/Fuj-tech-USB-Type-C-ulkoinen-kovalevykotelo-2-5-SATA-kovalev"}, "Fuj:tech USB Type-C -ulkoinen kovalevykotelo 2,5\" SATA-kovalevyille 31,90", false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
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
