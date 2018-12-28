package main

import (
	"context"
	"reflect"
	"testing"
)

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
		{"ValidURL1", args{url: "https://www.iltalehti.fi/kotimaa/a/bb594cd9-f66c-4bca-b626-a54848ea6ffb"}, "Poliisi kuvaa alakouluikäisen pojan murhaa erityisen raa’aksi, motiivi ja tekotapa vielä mysteeri - tämä kaikki tapauksesta tiedetään nyt", false},
		{"ValidURL2", args{url: "https://www.is.fi/taloussanomat/oma-raha/art-2000005935228.html"}, "Marko odottaa innolla talven sähkölaskuja – suosittu lämpöpumppu tuo jopa 1 500 euron säästön", false},
		{"ValidURL3", args{url: "https://yle.fi/uutiset/3-10507654"}, "Jätteiden mukana palaa miljoonien edestä arvokkaita metalleja – tutkijat löysivät menetelmän, jolla ne voidaan saada talteen tuhkasta", false},
		{"No opengraph", args{url: "https://www.manttavilppula.fi"}, "Taidekaupunki | Mänttä-Vilppula - Mänttä-Vilppula", false},
		{"No title", args{url: "http://addiktit.net"}, "", true},
		{"URL looks like jpg but isn't", args{url: "http://kuvaton.com/browse/57101/fatcop.jpg"}, "fatcop.jpg", false},
		{"Image", args{url: "https://i.imgur.com/r13Q6Yp.jpg"}, "", true},
		{"Does not exist", args{url: "https://definitely-not-a.website/"}, "", true},
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

func TestHandleRequest(t *testing.T) {
	type args struct {
		ctx   context.Context
		query TitleQuery
	}
	tests := []struct {
		name    string
		args    args
		want    TitleResponse
		wantErr bool
	}{
		{"ValidURL1", args{query: TitleQuery{URL: "https://www.iltalehti.fi/kotimaa/a/bb594cd9-f66c-4bca-b626-a54848ea6ffb"}}, TitleResponse{Title: "Poliisi kuvaa alakouluikäisen pojan murhaa erityisen raa’aksi, motiivi ja tekotapa vielä mysteeri - tämä kaikki tapauksesta tiedetään nyt"}, false},
		{"CustomYle", args{query: TitleQuery{URL: "https://areena.yle.fi/1-4192173"}}, TitleResponse{Title: "Areena custom handler"}, false},
		{"CustomApina", args{query: TitleQuery{URL: "https://apina.biz/167922"}}, TitleResponse{Title: ""}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HandleRequest(tt.args.ctx, tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandleRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HandleRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
