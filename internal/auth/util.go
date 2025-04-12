package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
)

func preExtractEmail(email string) (string, error) {
	split := strings.Split(email, "@")
	if len(split) != 2 {
		return "", fmt.Errorf("invalid email: %s", email)
	}

	return split[0], nil
}

func createSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}
