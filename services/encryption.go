package services

import (
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"

	"github.com/gofiber/utils"
	"golang.org/x/crypto/chacha20poly1305"
)

// EncryptionService - Interface for encryption and decryption
type EncryptionService interface {
	Encrypt(dst, msg []byte) ([]byte, error)
	EncryptString(msg string) ([]byte, error)
	Decrypt(msg []byte) (string, error)
}

type encryptionService struct {
	cipher cipher.AEAD
	key    []byte
}

// NewEncryptionService - Creates new instance of EncryptionService
func NewEncryptionService(key []byte) (EncryptionService, error) {
	c, err := chacha20poly1305.NewX(key)

	if err != nil {
		return nil, err
	}

	return encryptionService{
		key:    key,
		cipher: c,
	}, nil
}

func (s encryptionService) EncryptString(msg string) ([]byte, error) {
	capacity := s.cipher.NonceSize() + len(msg) + s.cipher.Overhead()

	dst := make([]byte, s.cipher.NonceSize(), capacity)

	return s.Encrypt(dst, utils.GetBytes(msg))
}

func (s encryptionService) Encrypt(dst, msg []byte) ([]byte, error) {
	capacity := s.cipher.NonceSize() + len(msg) + s.cipher.Overhead()

	if len(dst) != s.cipher.NonceSize() || cap(dst) != capacity {
		return nil, fmt.Errorf("Not enough bytes in dst, expected %d, given %d", capacity, cap(dst))
	}

	n, err := rand.Read(dst)

	if err != nil {
		return nil, err
	}

	if n != len(dst) {
		return nil, errors.New("Cannot generate random nonce")
	}

	return s.cipher.Seal(dst, dst, msg, nil), nil
}

func (s encryptionService) Decrypt(msg []byte) (string, error) {
	if len(msg) < s.cipher.NonceSize() {
		return "", errors.New("Size of message is less than nonce size")
	}
	nonce, ciphertext := msg[:s.cipher.NonceSize()], msg[s.cipher.NonceSize():]

	decrypted, err := s.cipher.Open(nil, nonce, ciphertext, nil)

	if err != nil {
		return "", err
	}

	return utils.GetString(decrypted), nil
}
