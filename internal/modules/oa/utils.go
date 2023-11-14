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

func BuildPaymentInformation(in *PayedAmount) *PaymentInformation {
	return &PaymentInformation{
		PayedAmount:     in,
		PaymentLocation: "PAY_LOCAL",
	}
}
