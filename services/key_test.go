package services

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestKeyPairGenerator(t *testing.T) {
	t.Parallel()
	asserts := require.New(t)

	t.Run("Generate", func(t *testing.T) {
		publicBuffer := bytes.NewBuffer(make([]byte, 0, 32))
		privateBuffer := bytes.NewBuffer(make([]byte, 0, 32))

		generator := NewKeyPairGenerator(publicBuffer, privateBuffer)
		asserts.Nilf(generator.Generate(), "Error while generating public and private key pair")
		asserts.Equal(publicBuffer.Len(), 32)
		asserts.Equal(privateBuffer.Len(), 32)
	})

	t.Run("GenerateWithOutput", func(t *testing.T) {
		publicKeyOut := make([]byte, 0, PublicKeyLength)
		privateKeyOut := make([]byte, 0, PrivateKeyLength)

		publicBuffer := bytes.NewBuffer(make([]byte, 0, PublicKeyLength))
		privateBuffer := bytes.NewBuffer(make([]byte, 0, PrivateKeyLength))

		generator := NewKeyPairGenerator(publicBuffer, privateBuffer)
		asserts.Nilf(generator.Generate(&publicKeyOut, &privateKeyOut), "Error while generating public and private key pair")
		asserts.Equal(publicBuffer.Len(), PublicKeyLength)
		asserts.Equal(privateBuffer.Len(), PrivateKeyLength)

		asserts.EqualValues(publicBuffer.Bytes(), publicKeyOut)
		asserts.EqualValues(privateBuffer.Bytes(), privateKeyOut)
	})
}

func TestSecretKeyGenerator(t *testing.T) {
	t.Parallel()
	asserts := require.New(t)
	publicBuffer := bytes.NewBuffer(make([]byte, 0, PublicKeyLength))
	privateBuffer := bytes.NewBuffer(make([]byte, 0, PrivateKeyLength))

	generator := NewKeyPairGenerator(publicBuffer, privateBuffer)
	asserts.Nil(generator.Generate())

	service, err := NewPublicKeyEncryption(publicBuffer, privateBuffer)
	asserts.Nil(err)

	t.Run("Generate", func(t *testing.T) {
		secretBuffer := bytes.NewBuffer(make([]byte, 0, SecretKeyLength))
		secretGenerator := NewSecretKeyGenerator(secretBuffer, service)
		asserts.Nil(secretGenerator.Generate())
	})

	t.Run("GenerateWithOutput", func(t *testing.T) {
		out := make([]byte, 0, SecretKeyLength)
		secretBuffer := bytes.NewBuffer(make([]byte, 0, SecretKeyLength))
		secretGenerator := NewSecretKeyGenerator(secretBuffer, service)
		asserts.Nil(secretGenerator.Generate(&out))

		decryptedSecret, err := service.Decrypt(nil, secretBuffer.Bytes())
		asserts.Nil(err)
		asserts.Len(decryptedSecret, SecretKeyLength)
		asserts.Len(out, SecretKeyLength)
		asserts.NotEqualValues(secretBuffer.Bytes(), out)
		asserts.EqualValues(decryptedSecret, out)
	})
}
