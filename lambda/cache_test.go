package lambda

import (
	"reflect"
	"testing"
)

func TestCheckCache(t *testing.T) {
	type args struct {
		query TitleQuery
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		//{"First test", args{TitleQuery{URL: "http://imgur.com/gallery/jqWjOLJ"}}, "any ideas how to clean this moldy leather interrior? [4 images]", false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := CheckCache(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckCache() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCacheAndReturn(t *testing.T) {
	type args struct {
		query TitleQuery
		title string
		err   error
	}
	tests := []struct {
		name    string
		args    args
		want    TitleQuery
		wantErr bool
	}{
		// TODO: Test cases with example.com urls
		// Test cleanup, delete example stuff?
		// {"First test", args{TitleQuery{URL: "http://imgur.com/gallery/jqWjOLJ"}}, "any ideas how to clean this moldy leather interrior? [4 images]", false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := CacheAndReturn(tt.args.query, tt.args.title, tt.args.err)
			if (err != nil) != tt.wantErr {
				t.Errorf("CacheAndReturn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CacheAndReturn() = %v, want %v", got, tt.want)
			}
		})
	}
}
