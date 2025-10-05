package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// Создает HMAC-SHA256 хеш от строки с секретным ключом
func CreateHash(data string, secretKey string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(data))

	return hex.EncodeToString(h.Sum(nil))
}

// Проверяет валидность HMAC хеша
func VerifyHash(data string, hash string, secretKey string) bool {
	expectedHash := CreateHash(data, secretKey)
	return hmac.Equal([]byte(expectedHash), []byte(hash))
}
