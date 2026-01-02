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
	key, err := GenerateMasterKey()
	if err != nil {
		t.Fatal(err)
	}

	if len(key) != 32 {
		t.Error("AES-256 master key is not 256 bits long")
	}
}
