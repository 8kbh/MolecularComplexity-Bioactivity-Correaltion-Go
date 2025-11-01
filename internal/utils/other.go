package utils

import (
	"crypto/rand"
	"math/big"
)

// GenerateRandomString generates a random string of the specified length.
// The string contains uppercase letters, lowercase letters, and digits.
func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)

	for i := 0; i < length; i++ {
		// Generate a random index within the charset length
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return ""
		}
		result[i] = charset[num.Int64()]
	}

	return string(result)
}
