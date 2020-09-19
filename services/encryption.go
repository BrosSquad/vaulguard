package services

import (
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/gofiber/utils"
	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/nacl/box"
	"io"
)

var (
	ErrNotEnoughBytes = errors.New("not enough bytes read from crypto random source")
	ErrKeyLength      = errors.New("key has to be 32 bytes long")
)

type SecretKeyGeneration interface {
	GenerateEncryptionKey() ([]byte, error)
}

type PublicKeyCryptographyKeysGeneration interface {
	GenerateKeyPairs() (*[32]byte, *[32]byte, error)
}

// EncryptionService - Interface for encryption and decryption
type EncryptionService interface {
	Encrypt(dst, msg []byte) ([]byte, error)
	EncryptString(msg string) ([]byte, error)
	Decrypt(dst, msg []byte) ([]byte, error)
	DecryptString(msg []byte) (string, error)
}

type secretKeyEncryption struct {
	cipher cipher.AEAD
	key    []byte
}

func (s secretKeyEncryption) GenerateEncryptionKey() ([]byte, error) {
	key := make([]byte, chacha20poly1305.KeySize)
	n, err := rand.Read(key)

	if err != nil {
		return nil, err
	}

	if n != chacha20poly1305.KeySize {
		return nil, ErrNotEnoughBytes
	}

	return key, nil
}

// NewSecretKeyEncryption - Creates new instance of EncryptionService
func NewSecretKeyEncryption(key []byte) (EncryptionService, error) {
	c, err := chacha20poly1305.NewX(key)

	if err != nil {
		return nil, err
	}

	return secretKeyEncryption{
		key:    key,
		cipher: c,
	}, nil
}

func (s secretKeyEncryption) EncryptString(msg string) ([]byte, error) {
	capacity := s.cipher.NonceSize() + len(msg) + s.cipher.Overhead()

	dst := make([]byte, s.cipher.NonceSize(), capacity)

	return s.Encrypt(dst, utils.GetBytes(msg))
}

func (s secretKeyEncryption) Encrypt(dst, msg []byte) ([]byte, error) {
	capacity := s.cipher.NonceSize() + len(msg) + s.cipher.Overhead()

	if len(dst) != s.cipher.NonceSize() || cap(dst) != capacity {
		return nil, fmt.Errorf("not enough bytes in dst, expected %d, given %d", capacity, cap(dst))
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

func (s secretKeyEncryption) Decrypt(dst, msg []byte) ([]byte, error) {
	if len(msg) < s.cipher.NonceSize() {
		return nil, errors.New("size of message is less than nonce size")
	}
	nonce, ciphertext := msg[:s.cipher.NonceSize()], msg[s.cipher.NonceSize():]

	decrypted, err := s.cipher.Open(dst, nonce, ciphertext, nil)

	if err != nil {
		return nil, err
	}

	return decrypted, nil
}

func (s secretKeyEncryption) DecryptString(msg []byte) (string, error) {
	message, err := s.Decrypt(nil, msg)

	if err != nil {
		return "", err
	}

	return utils.GetString(message), nil
}

func readKey(out []byte, r io.Reader) error {
	read, err := r.Read(out)
	if err != nil {
		return err
	}

	if read != 32 {
		return ErrKeyLength
	}

	return nil
}

func NewPublicKeyEncryption(publicKey, privateKey io.Reader) (EncryptionService, error) {
	var publicKeyBytes [32]byte
	var privateKeyBytes [32]byte

	if err := readKey(publicKeyBytes[:], publicKey); err != nil {
		return nil, err
	}

	if err := readKey(privateKeyBytes[:], privateKey); err != nil {
		return nil, err
	}

	return publicKeyEncryption{
		privateKey: &privateKeyBytes,
		publicKey:  &publicKeyBytes,
	}, nil
}

type publicKeyEncryption struct {
	privateKey *[32]byte
	publicKey  *[32]byte
}

func (p publicKeyEncryption) GenerateKeyPairs() (*[32]byte, *[32]byte, error) {
	return box.GenerateKey(rand.Reader)
}

func (p publicKeyEncryption) Encrypt(dst, msg []byte) ([]byte, error) {
	return box.SealAnonymous(dst, msg, p.publicKey, rand.Reader)
}

func (p publicKeyEncryption) EncryptString(msg string) ([]byte, error) {
	dst := make([]byte, 0)
	return p.Encrypt(dst, utils.GetBytes(msg))
}

func (p publicKeyEncryption) Decrypt(dst, msg []byte) ([]byte, error) {
	message, ok := box.OpenAnonymous(dst, msg, p.publicKey, p.privateKey)

	if !ok {
		return nil, errors.New("decryption failed")
	}

	return message, nil
}

func (p publicKeyEncryption) DecryptString(msg []byte) (string, error) {
	message, err := p.Decrypt(nil, msg)

	if err != nil {
		return "", err
	}

	return utils.GetString(message), nil
}
