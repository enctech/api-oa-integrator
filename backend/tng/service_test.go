package tng

import (
	"reflect"
	"testing"
)

func Test_calculatePercentSurchargeAmount(t *testing.T) {
	type args struct {
		txn       float64
		tax       float64
		surcharge float64
	}
	tests := []struct {
		name string
		args args
		want TaxCalculation
	}{
		{
			name: "Test based on defined value",
			args: args{txn: 5, surcharge: 6, tax: 5},
			want: TaxCalculation{
				surcharge:       0.22,
				surchargeAmt:    0.27,
				surchargeTaxAmt: 0.02,
				parkingAmt:      4.49,
				parkingTaxAmt:   0.22,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculatePercentSurchargeAmount(tt.args.txn, tt.args.tax, tt.args.surcharge); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("calculatePercentSurchargeAmount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_calculateExactSurchargeAmount(t *testing.T) {
	type args struct {
		txn       float64
		tax       float64
		surcharge float64
	}
	tests := []struct {
		name string
		args args
		want TaxCalculation
	}{
		{
			name: "Test based on defined value",
			args: args{txn: 5, tax: 5, surcharge: 0.10},
			want: TaxCalculation{
				surcharge:       0.10,
				surchargeAmt:    0.10,
				surchargeTaxAmt: 0.00,
				parkingAmt:      4.67,
				parkingTaxAmt:   0.23,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateExactSurchargeAmount(tt.args.txn, tt.args.tax, tt.args.surcharge); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("calculateExactSurchargeAmount() = %v, want %v", got, tt.want)
			}
		})
	}
}
