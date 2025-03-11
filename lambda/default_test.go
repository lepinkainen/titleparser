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
		{"No opengraph", args{url: "https://www.manttavilppula.fi"}, "Mänttä-Vilppula | Taidekaupunki keskellä kaunista järvimaisemaa", false},
		{"URL looks like jpg but isn't", args{url: "http://kuvaton.com/browse/57101/fatcop.jpg"}, "fatcop.jpg", false},
		{"Image", args{url: "https://i.imgur.com/r13Q6Yp.jpg"}, "", true},
		{"Bleepingcomputer 403", args{url: "https://www.bleepingcomputer.com/news/apple/apple-fixes-webkit-zero-day-exploited-in-extremely-sophisticated-attacks/"}, "", true},
		{"Yahoo finance 429", args{url: "https://finance.yahoo.com/news/tesla-stock-tumbles-over-15-wiping-out-post-election-gains-as-demand-worries-continue-to-weigh-161337256.html"}, "", true},
		{"Gogdotcom", args{url: "https://www.gog.com/game/diablo"}, "Diablo + Hellfire", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DefaultHandler(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("DefaultHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DefaultHandler() = '%v', want '%v'", got, tt.want)
			}
		})
	}
}
