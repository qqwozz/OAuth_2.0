package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GenerateRandomString() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("crypto/rand failed: %v", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
