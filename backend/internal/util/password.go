package util

import (
	"crypto/rand"
	"math/big"
)

func GeneratePassword(length int) (string, error) {
	// 密码字符集
	const charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"0123456789" +
		"!@#$%^&*()_+"

	var password = make([]byte, length)
	for i := 0; i < length; i++ {
		charIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		password[i] = charset[charIndex.Int64()]
	}

	return string(password), nil
}
