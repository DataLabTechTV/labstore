package security

import (
	"crypto/rand"
	"math/big"
)

const passwordCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Produces an alphanumeric placeholder password
func GeneratePassword(length int) (string, error) {
	randBytes := make([]byte, length)
	for i := range randBytes {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(passwordCharset))))
		if err != nil {
			return "", err
		}
		randBytes[i] = passwordCharset[num.Int64()]
	}

	return string(randBytes), nil
}

// Produces a key with the given byte-length
func GenerateKey(length int) ([]byte, error) {
	key := make([]byte, length)

	if _, err := rand.Read(key); err != nil {
		return nil, err
	}

	return key, nil
}

// Produces a AES-256 master key (32 bytes)
func GenerateMasterKey() ([]byte, error) {
	return GenerateKey(32)
}
