package services

import (
	"crypto/rand"
	"testing"
)

func TestEncryptionService(t *testing.T) {
	key := make([]byte, 32)

	_, err := rand.Read(key)

	if err != nil {
		t.Errorf("Cannot generate random key: %v", err)
	}

	service, err := NewEncryptionService(key)

	if err != nil {
		t.Errorf("Cannot create encryption service: %v", err)
	}

	t.Run("Encryption", func(t *testing.T) {
		encryptedBytes, err := service.EncryptString("Hello World")

		if err != nil {
			t.Errorf("Cannot encrypt text: %v", err)
		}

		str, err := service.Decrypt(encryptedBytes)

		if err != nil {
			t.Errorf("Decryption failed: %v", err)
		}

		if str != "Hello World" {
			t.Errorf("Starting string is not equal to decrypted string")
		}
	})

	t.Run("SmallMessageInDecryption", func(t *testing.T) {
		dst := make([]byte, 24, 25)

		_, err := service.Encrypt(dst, []byte("Hello World"))

		if err == nil {
			t.Error(err)
		}
	})

	t.Run("DecryptionWithSmallMessageSize", func(t *testing.T) {
		data := make([]byte, 12)

		_, _ = rand.Read(data)

		_, err := service.Decrypt(data)

		if err == nil {
			t.Error(err)
		}
	})
}
