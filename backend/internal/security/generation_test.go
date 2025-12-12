package security

import "testing"

const passwordLength = 32

func TestGeneratePassword(t *testing.T) {
	password, err := GeneratePassword(passwordLength)
	if err != nil {
		t.Fatal(err)
	}

	if len(password) != passwordLength {
		t.Error("Password does not match requested length")
	}
}

func TestGenerateAES256MasterKey(t *testing.T) {
	key, err := GenerateAES256MasterKey()
	if err != nil {
		t.Fatal(err)
	}

	if len(key) != 32 {
		t.Error("AES-256 master key is not 256 bits long")
	}
}

func TestGenerateSalt(t *testing.T) {
	key, err := GenerateSalt()
	if err != nil {
		t.Fatal(err)
	}

	if len(key) != 16 {
		t.Error("Salt is not 128 bits long")
	}
}
