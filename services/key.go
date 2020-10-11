package services

import (
	"crypto/rand"
	"errors"
	"golang.org/x/crypto/nacl/box"
	"io"
)

type keyPairGenerator struct {
	public  io.Writer
	private io.Writer
}

type secretKeyGenerator struct {
	encryption Encryption
	out        io.Writer
}

type Generator interface {
	Generate(out ...interface{}) error
}

func NewKeyPairGenerator(public, private io.Writer) Generator {
	return keyPairGenerator{public, private}
}

func NewSecretKeyGenerator(out io.Writer, service Encryption) Generator {
	return secretKeyGenerator{
		encryption: service,
		out:        out,
	}
}

func (s secretKeyGenerator) Generate(out ...interface{}) error {
	key := make([]byte, SecretKeyLength)
	n, err := rand.Read(key)

	if err != nil {
		return err
	}

	if n != SecretKeyLength {
		return ErrNotEnoughBytes
	}

	if len(out) == 1 {
		setOut(out[0], key)
	}

	enc, err := s.encryption.Encrypt(nil, key)

	if err != nil {
		return err
	}

	return writeKey(s.out, enc)
}

func (g keyPairGenerator) Generate(out ...interface{}) error {
	public, private, err := box.GenerateKey(rand.Reader)

	if err != nil {
		return err
	}

	if err := writeKey(g.public, public[:]); err != nil {
		return err
	}

	if err := writeKey(g.private, private[:]); err != nil {
		return err
	}

	if len(out) == 2 {
		setOut(out[0], public[:])
		setOut(out[1], private[:])
	}

	return nil
}

func writeKey(w io.Writer, key []byte) error {
	n, err := w.Write(key[:])

	if err != nil {
		return err
	}

	if n != len(key) {
		return errors.New("not enough bytes written to io.Writer")
	}

	return nil
}

func setOut(out interface{}, value []byte) {
	bytes := out.(*[]byte)
	*bytes = value
}
