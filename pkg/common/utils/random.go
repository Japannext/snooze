package utils

import (
	"crypto/rand"
	"encoding/base64"
)

// Return a random base64 string representing `n` bytes
func RandomURLSafeBase64(n int) (string, error) {
	b := make([]byte, 96)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
