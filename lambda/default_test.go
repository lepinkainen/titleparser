package lambda

import (
	"testing"
)

func TestDefaultHandler(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"ValidURL2", args{url: "https://www.is.fi/taloussanomat/oma-raha/art-2000005935228.html"}, "Marko odottaa innolla talven sähkölaskuja – suosittu lämpöpumppu tuo jopa 1 500 euron säästön", false},
		{"ValidURL3", args{url: "https://yle.fi/uutiset/3-10507654"}, "Jätteiden mukana palaa miljoonien edestä arvokkaita metalleja – tutkijat löysivät menetelmän, jolla ne voidaan saada talteen tuhkasta", false},
		{"No opengraph", args{url: "https://www.manttavilppula.fi"}, "Taidekaupunki | Mänttä-Vilppula - Mänttä-Vilppula", false},
		//{"No title", args{url: "http://addiktit.net"}, "", true},
		{"URL looks like jpg but isn't", args{url: "http://kuvaton.com/browse/57101/fatcop.jpg"}, "fatcop.jpg", false},
		{"Image", args{url: "https://i.imgur.com/r13Q6Yp.jpg"}, "", true},
		//TODO: bell-labs added an og:title, which isn't broken find another example
		//{"Whitespace in middle of title", args{url: "https://www.bell-labs.com/unix50/"}, "Unix 50", false},
		// TODO: Gog.com localises based on IP, so this fails in multiple ways depending on where the build machine is...
		//{"Gog.com Diablo", args{url: "https://www.gog.com/game/diablo"}, "Diablo + Hellfire on GOG.com", false},
		// TODO: Long title (over 200 characters)
		//{"Does not exist", args{url: "https://definitely-not-a.website/"}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DefaultHandler(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("DefaultHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DefaultHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
