package main

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateKeyPair(t *testing.T) {
	t.Parallel()
	asserts := require.New(t)

	path, err := ioutil.TempDir("", "generate_vaulguard_keys")
	asserts.Nil(err)

	t.Run("SuccessfulGeneration", func(t *testing.T) {
		publicKeyPath := filepath.Join(path, "public_success.key")
		privateKeyPath := filepath.Join(path, "private_success.key")
		_, err := generateKeyPair(privateKeyPath, publicKeyPath, true)
		asserts.Nil(err)
		asserts.FileExists(publicKeyPath)
		asserts.FileExists(privateKeyPath)
		asserts.Nil(os.Remove(publicKeyPath))
		asserts.Nil(os.Remove(privateKeyPath))
	})

	t.Run("PrivateKeyNotFound", func(t *testing.T) {
		publicKeyPath := filepath.Join(path, "public_success.key")
		privateKeyPath := filepath.Join(path, "private_success.key")
		file, err := os.Create(publicKeyPath)
		asserts.Nil(err)
		asserts.Nil(file.Close())
		asserts.FileExists(publicKeyPath)
		_, err = generateKeyPair(privateKeyPath, publicKeyPath, false)
		asserts.NotNil(err)
	})

	t.Run("PublicKeyNotFound", func(t *testing.T) {
		publicKeyPath := filepath.Join(path, "public_success.key")
		privateKeyPath := filepath.Join(path, "private_success.key")
		file, err := os.Create(privateKeyPath)
		asserts.Nil(err)
		asserts.Nil(file.Close())
		asserts.FileExists(publicKeyPath)
		_, err = generateKeyPair(privateKeyPath, publicKeyPath, false)
		asserts.NotNil(err)
	})

	t.Run("EmptyKeyFiles", func(t *testing.T) {
		publicKeyPath := filepath.Join(path, "public_success.key")
		privateKeyPath := filepath.Join(path, "private_success.key")
		file1, err := os.Create(publicKeyPath)
		asserts.Nil(err)
		asserts.Nil(file1.Close())
		file2, err := os.Create(privateKeyPath)
		asserts.Nil(err)
		asserts.Nil(file2.Close())
		asserts.FileExists(publicKeyPath)
		asserts.FileExists(privateKeyPath)
		_, err = generateKeyPair(privateKeyPath, publicKeyPath, false)
		asserts.NotNil(err)
	})

}
