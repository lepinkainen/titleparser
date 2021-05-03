package handler

import "testing"

func TestHumanizeNumber(t *testing.T) {
	type args struct {
		views int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"ones", args{views: 10}, "10"},
		{"thousands", args{views: 1000}, "1k"},
		{"millions", args{views: 1000000}, "1M"},
		{"billions", args{views: 1000000000}, "1Billion"},
		{"trillions", args{views: 1000000000000}, "1Trillion"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := HumanizeNumber(tt.args.views); got != tt.want {
				t.Errorf("HumanizeNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}
