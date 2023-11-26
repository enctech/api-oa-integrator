package oa

import (
	"fmt"
	"hash/crc32"
)

func encryptLpn(in string) string {
	// Use CRC32 for simplicity
	crc := crc32.ChecksumIEEE([]byte(in))

	// Take the modulus to get an 8-digit number
	encrypted := crc % 100000000

	// Format as an 8-digit number
	return fmt.Sprintf("%08d", encrypted)
}

func BuildPaymentInformation(in *PaymentData) *PaymentInformation {
	if in == nil {
		return &PaymentInformation{
			PaymentLocation: "PAY_LOCAL",
		}
	}
	return &PaymentInformation{
		PayedAmount: &PayedAmount{
			Amount:  in.OriginalAmount.Amount,
			VatRate: in.OriginalAmount.VatRate,
		},
		PaymentLocation: "PAY_LOCAL",
	}
}
