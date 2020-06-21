package handler

import "testing"

func TestPr0gramm(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"test1", args{url: "https://pr0gramm.com/top/Pussy%20Massage/3974894"}, "", false},
		{"test2", args{url: "https://pr0gramm.com/top/Triggerhochlad/3974882"}, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Pr0gramm(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("Pr0gramm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Pr0gramm() = %v, want %v", got, tt.want)
			}
		})
	}
}
