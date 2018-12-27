package main

import "testing"

func TestFindTitle(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"URL1", args{url: "https://www.iltalehti.fi/kotimaa/a/bb594cd9-f66c-4bca-b626-a54848ea6ffb"}, "Poliisi kuvaa alakouluikäisen pojan murhaa erityisen raa’aksi, motiivi ja tekotapa vielä mysteeri - tämä kaikki tapauksesta tiedetään nyt", false},
		{"URL2", args{url: "https://www.is.fi/taloussanomat/oma-raha/art-2000005935228.html"}, "Marko odottaa innolla talven sähkölaskuja – suosittu lämpöpumppu tuo jopa 1 500 euron säästön", false},
		{"No opengraph", args{url: "https://www.manttavilppula.fi"}, "Taidekaupunki | Mänttä-Vilppula - Mänttä-Vilppula", false},
		{"No title", args{url: "http://addiktit.net"}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindTitle(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindTitle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FindTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}
