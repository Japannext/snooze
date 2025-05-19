package utils

import (
	"crypto/rand"
	"encoding/base64"

	log "github.com/sirupsen/logrus"
)

// Return a random base64 string representing `n` bytes.
func RandomURLSafeBase64(n int) string {
	b := make([]byte, n)

	_, err := rand.Read(b)
	if err != nil {
		log.Fatalf("failed to read random: %s", err)
	}

	return base64.URLEncoding.EncodeToString(b)
}
