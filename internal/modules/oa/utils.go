package oa

import (
	"fmt"
	"strings"
)

func encrypt(in string) string {
	encoded := fmt.Sprintf("%08x", in)

	// Ensure it's exactly 8 digits by truncating or padding
	encoded = encoded[:8]

	return strings.ToUpper(encoded)
}
