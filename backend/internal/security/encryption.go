package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"log/slog"
	"os"
)

func NewAESGCM(masterKeyPath string) (cipher.AEAD, error) {
	key, err := os.ReadFile(masterKeyPath)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return aesGCM, nil
}

// Encrypt plain text using AES-GCM, returning nonce+ciphertext
func EncryptAESGCM(plainText string, masterKeyPath string) ([]byte, error) {
	aesGCM, err := NewAESGCM(masterKeyPath)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := aesGCM.Seal(nil, nonce, []byte(plainText), nil)

	encrypted := make([]byte, 0, len(nonce)+len(ciphertext))
	encrypted = append(encrypted, nonce...)
	encrypted = append(encrypted, ciphertext...)

	slog.Debug(
		"encrypt aes gcm",
		"len(nonce)", len(nonce),
		"len(cipherText)", len(ciphertext),
		"len(encrypted)", len(encrypted),
	)

	return encrypted, nil
}

// Decrypt data encrypted with AES-GCM (nonce+ciphertext), returning the plain text
func DecryptAESGCM(encrypted []byte, masterKeyPath string) (string, error) {
	aesGCM, err := NewAESGCM(masterKeyPath)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	nonce := encrypted[:nonceSize]
	ciphertext := encrypted[nonceSize:]

	slog.Debug(
		"decrypt aes gcm",
		"len(encrypted)", len(encrypted),
		"len(nonce)", len(nonce),
		"len(cipherText)", len(ciphertext),
	)

	decrypted, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}
