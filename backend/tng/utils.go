package tng

import (
	"crypto/sha256"
	"encoding/base64"
)

func createSignature(in string) string {
	hasher := sha256.New()
	hasher.Write([]byte(in))
	hashBytes := hasher.Sum(nil)

	// Encode the hash in base64
	hashBase64 := base64.StdEncoding.EncodeToString(hashBytes)
	return hashBase64
}
