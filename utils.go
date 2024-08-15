package main

import (
	"crypto/rand"
	"encoding/base32"
)

func generateNanoID() (string, error) {
	random_bytes := make([]byte, 15)
	_, err := rand.Read(random_bytes)
	if err != nil {
		return "", err
	}
	return base32.HexEncoding.EncodeToString(random_bytes), nil
}
