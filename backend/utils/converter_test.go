package utils

import "testing"

func Test_roundMoney(t *testing.T) {
	type args struct {
		amount float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "Test value 1",
			args: args{amount: 10.995},
			want: 11.00,
		},
		{
			name: "Test value 2",
			args: args{amount: 10.456},
			want: 10.46,
		},
		{
			name: "Test value 3",
			args: args{amount: 10.454},
			want: 10.45,
		},
		{
			name: "Test value 4",
			args: args{amount: 6.782},
			want: 6.78,
		},
		{
			name: "Test value 5",
			args: args{amount: 6.787},
			want: 6.79,
		},
		{
			name: "Test value 6",
			args: args{amount: 6.785876950},
			want: 6.79,
		},
		{
			name: "Test value 7",
			args: args{amount: 6.7800000001},
			want: 6.78,
		},
		{
			name: "Test value 8",
			args: args{amount: 6.7800005001},
			want: 6.78,
		},
		{
			name: "Test value 9",
			args: args{amount: 6.785},
			want: 6.79,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RoundMoney(tt.args.amount); got != tt.want {
				t.Errorf("roundMoney() = %v, want %v", got, tt.want)
			}
		})
	}
}
