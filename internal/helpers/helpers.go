package helpers

import (
	"cafeteria/internal/models"
	"crypto/md5"
	"encoding/hex"
	"io"
)

func IsValid(t models.TransactionType) bool {
	switch t {
	case 0, 1, 2:
		return true
	default:
		return false
	}
}

func CreateMd5Hash(text string) string {
	hasher := md5.New()
	_, err := io.WriteString(hasher, text)
	if err != nil {
		// panic(err)
	}

	return hex.EncodeToString(hasher.Sum(nil))
}
