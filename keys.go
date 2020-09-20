package main

import (
	"errors"
	"io"
	"io/ioutil"
	"os"

	"github.com/BrosSquad/vaulguard/config"
	"github.com/BrosSquad/vaulguard/services"
)

const DefaultPermission = 0700

func generateKeyPair(privateKeyPath, publicKeyPath string, create bool) (services.Encryption, error) {
	if err := createDirs(privateKeyPath, publicKeyPath); err != nil {
		return nil, err
	}
	flags := os.O_RDWR | os.O_CREATE

	privateKeyFile, err := os.OpenFile(privateKeyPath, flags, DefaultPermission)
	if err != nil {
		return nil, err
	}
	defer privateKeyFile.Close()

	publicKeyFile, err := os.OpenFile(publicKeyPath, flags, DefaultPermission)
	if err != nil {
		return nil, err
	}
	defer publicKeyFile.Close()
	if create {
		keyGenerator := services.NewKeyPairGenerator(publicKeyFile, privateKeyFile)
		// Generate key pair
		if err := keyGenerator.Generate(); err != nil {
			return nil, err
		}
		if _, err = publicKeyFile.Seek(0, io.SeekStart); err != nil {
			return nil, err
		}
		if _, err = privateKeyFile.Seek(0, io.SeekStart); err != nil {
			return nil, err
		}
	}
	return services.NewPublicKeyEncryption(publicKeyFile, privateKeyFile)
}

func getSecretKey(service services.Encryption, secretKeyPath string, secretKeyExists bool) ([]byte, error) {
	if err := createDirs(secretKeyPath); err != nil {
		return nil, err
	}

	flags := os.O_RDWR | os.O_CREATE
	secretKeyFile, err := os.OpenFile(secretKeyPath, flags, DefaultPermission)
	if err != nil {
		return nil, err
	}
	defer secretKeyFile.Close()
	if secretKeyExists {
		key, err := ioutil.ReadAll(secretKeyFile)

		if err != nil {
			return nil, err
		}

		return service.Decrypt(nil, key)
	}

	var key []byte
	if err := services.NewSecretKeyGenerator(secretKeyFile, service).Generate(&key); err != nil {
		return nil, err
	}

	return key, nil
}

func checkKeyPairExistence(privateKeyExists, publicKeyExists bool) error {
	if privateKeyExists && !publicKeyExists {
		return errors.New("public key does not exit while private exists")
	}

	if !privateKeyExists && publicKeyExists {
		return errors.New("private key does not exit while public exists")
	}

	return nil
}

func getKeys(config *config.Config) ([]byte, error) {
	privateKeyPath, err := getAbsolutePath(config.Keys.Private)

	if err != nil {
		return nil, err
	}

	publicKeyPath, err := getAbsolutePath(config.Keys.Public)
	if err != nil {
		return nil, err
	}

	secretKeyPath, err := getAbsolutePath(config.Keys.Secret)
	if err != nil {
		return nil, err
	}

	privateKeyExists, publicKeyExists, secretKeyExists := fileExists(privateKeyPath), fileExists(publicKeyPath), fileExists(secretKeyPath)

	if err := checkKeyPairExistence(privateKeyExists, publicKeyExists); err != nil {
		return nil, err
	}

	service, err := generateKeyPair(privateKeyPath, publicKeyPath, !publicKeyExists)
	if err != nil {
		return nil, err
	}

	return getSecretKey(service, secretKeyPath, secretKeyExists)
}
