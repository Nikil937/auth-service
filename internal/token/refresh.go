package token

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func GenerateRefreshToken() (string, error) {
	bytes := make([]byte, 32)

	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generate random bytes: %w", err)
	}

	return hex.EncodeToString(bytes), nil
}
