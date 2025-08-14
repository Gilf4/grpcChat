package refresh

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

const (
	lengthChars = 512

	bytesNeeded = 384
)

func GenerateToken() (string, error) {
	op := "lib.refresh.GenerateToken"

	b := make([]byte, bytesNeeded)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	token := base64.RawURLEncoding.EncodeToString(b)

	return token, nil
}
