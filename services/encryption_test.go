package services

import (
	"bytes"
	"crypto/rand"
	"golang.org/x/crypto/nacl/box"
	"testing"
)

func TestSecretKeyService(t *testing.T) {
	t.Parallel()
	key := make([]byte, 32)

	_, err := rand.Read(key)

	if err != nil {
		t.Fatalf("Cannot generate random key: %v\n", err)
	}

	service, err := NewSecretKeyEncryption(key)

	if err != nil {
		t.Fatalf("Cannot create encryption service: %v\n", err)
	}

	t.Run("Encryption", func(t *testing.T) {
		encryptedBytes, err := service.EncryptString("Hello World")

		if err != nil {
			t.Fatalf("Cannot encrypt text: %v", err)
		}

		str, err := service.DecryptString(encryptedBytes)

		if err != nil {
			t.Fatalf("Decryption failed: %v", err)
		}

		if str != "Hello World" {
			t.Fatalf("Starting string is not equal to decrypted string\n")
		}
	})

	t.Run("SmallMessageInDecryption", func(t *testing.T) {
		dst := make([]byte, 24, 25)

		_, err := service.Encrypt(dst, []byte("Hello World"))

		if err == nil {
			t.Fatal(err)
		}
	})

	t.Run("DecryptionWithSmallMessageSize", func(t *testing.T) {
		data := make([]byte, 12)

		_, _ = rand.Read(data)

		_, err := service.DecryptString(data)

		if err == nil {
			t.Fatal(err)
		}
	})
}

func TestPublicKeyService(t *testing.T) {
	t.Parallel()
	public, private, err := box.GenerateKey(rand.Reader)

	if err != nil {
		t.Fatalf("Error while generating public and private key pair: %v\n", err)
	}

	publicKeyBuf := bytes.NewBuffer(public[:])
	privateKeyBuf := bytes.NewBuffer(private[:])

	service, err := NewPublicKeyEncryption(publicKeyBuf, privateKeyBuf)

	if err != nil {
		t.Fatalf("Error while creating public key encryption service: %v\n", err)
	}

	t.Run("Encrypt", func(t *testing.T) {
		data, err := service.EncryptString("Hello World")

		if err != nil {
			t.Fatalf("Error while encrypting: %v\n", err)
		}

		if data == nil || len(data) == 0 {
			t.Fatalf("Error while encrypting, no encrypted data\n")
		}
	})

	t.Run("Decryption", func(t *testing.T) {
		data, err := service.EncryptString("Hello World")

		if err != nil {
			t.Fatalf("Error while encrypting: %v\n", err)
		}

		if data == nil || len(data) == 0 {
			t.Fatalf("Error while encrypting, no encrypted data\n")
		}

		message, err := service.DecryptString(data)

		if err != nil {
			t.Fatalf("Error while decrypting: %v\n", err)
		}

		if message != "Hello World" {
			t.Fatalf("Error while decrypting: Expected message \"Hello World\", Given: %s", message)
		}
	})

}
