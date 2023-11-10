package tng

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
)

func createSignature(filePath string) string {
	getwd, err := os.Getwd()
	if err != nil {
		return ""
	}
	fmt.Println(getwd)
	fileData, err := os.ReadFile(fmt.Sprintf("%s/%s", getwd, filePath))
	if err != nil {
		fmt.Println("Error reading file:", err)
		return ""
	}

	// Calculate the SHA-256 hash
	hasher := sha256.New()
	hasher.Write(fileData)
	hashBytes := hasher.Sum(nil)

	// Encode the hash in base64
	hashBase64 := base64.StdEncoding.EncodeToString(hashBytes)
	return hashBase64
}
