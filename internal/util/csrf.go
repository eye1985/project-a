package util

import (
	"encoding/base64"
	"log"
	"os"
)

func GetCSRFKey() []byte {
	base64Key := os.Getenv("CSRF_AUTH_KEY")
	key, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		log.Fatalf("Invalid base64 CSRF key: %v", err)
	}
	if len(key) != 32 {
		log.Fatalf("CSRF key must be 32 bytes, got %d", len(key))
	}
	return key
}
