package tng

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
)

// Signer provides digital signature functionality
type Signer struct {
	privateKey *rsa.PrivateKey
}

// NewSigner creates a new Signer with the provided PEM-encoded private key
func NewSigner(privateKeyPem string) (*Signer, error) {
	privateKey, err := parsePrivateKey(privateKeyPem)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return &Signer{
		privateKey: privateKey,
	}, nil
}

// parsePrivateKey parses a PEM-encoded private key
func parsePrivateKey(privateKeyPem string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privateKeyPem))
	if block == nil {
		return nil, errors.New("failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
	}

	rsaKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("not an RSA private key")
	}

	return rsaKey, nil
}

func (s *Signer) Sign(data string) (string, error) {
	dataTrimmed := strings.TrimSpace(data)
	hashed := sha256.Sum256([]byte(dataTrimmed))

	signature, err := rsa.SignPKCS1v15(rand.Reader, s.privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return "", fmt.Errorf("failed to sign data: %w", err)
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}
