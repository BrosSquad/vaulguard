package services

import (
	"crypto/rand"
	"testing"
)

func setupEncryption(t *testing.T) EncryptionService {
	key := make([]byte, 32)

	_, err := rand.Read(key)

	if err != nil {
		t.Errorf("Cannot generate random key: %v", err)
	}

	service, err := NewEncryptionService(key)

	if err != nil {
		t.Errorf("Cannot create encryption service: %v", err)
	}

	return service
}

func TestEncryptionDecryption(t *testing.T) {
	service := setupEncryption(t)

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
}

func TestEncryptionWithSmallDestination(t *testing.T) {
	service := setupEncryption(t)

	dst := make([]byte, 24, 25)

	_, err := service.Encrypt(dst, []byte("Hello World"))

	if err == nil {
		t.Error(err)
	}
}
func TestDecryptionnWithSmallMessageSize(t *testing.T) {
	service := setupEncryption(t)

	data := make([]byte, 12)

	rand.Read(data)

	_, err := service.Decrypt(data)

	if err == nil {
		t.Error(err)
	}
}
