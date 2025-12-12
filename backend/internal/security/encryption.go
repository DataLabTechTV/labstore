package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"os"
)

type EncryptedData struct {
	Value []byte
	Salt  []byte
}

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

// Encrypt plain text using AES-GCM, returning the ciphered text alongside the salt
func EncryptAESGCM(plainText string, masterKeyPath string) (*EncryptedData, error) {
	aesGCM, err := NewAESGCM(masterKeyPath)
	if err != nil {
		return nil, err
	}

	salt := make([]byte, aesGCM.NonceSize())
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	cipheredData := aesGCM.Seal(nil, salt, []byte(plainText), nil)

	encrypted := &EncryptedData{
		Value: cipheredData,
		Salt:  salt,
	}

	return encrypted, nil
}

// Decrypt data encrypted with AES-GCM, returning the plain text
func DecryptAESGCM(encrypted *EncryptedData, masterKeyPath string) (string, error) {
	aesGCM, err := NewAESGCM(masterKeyPath)
	if err != nil {
		return "", err
	}

	decrypted, err := aesGCM.Open(nil, encrypted.Salt, encrypted.Value, nil)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}
