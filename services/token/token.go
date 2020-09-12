package token

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"github.com/BrosSquad/vaulguard/models"
	"lukechampine.com/blake3"
)

type Service interface {
	Generate(uint) string
	Verify(string) (models.Application, bool)
}

type service struct {
	storage Storage
}

func NewService(storage Storage) Service {
	return service{storage}
}

func (s service) Generate(applicationId uint) string {
	tokenBytes := make([]byte, 64)
	_, err := rand.Read(tokenBytes)

	if err != nil {
		return ""
	}

	hashed := blake3.Sum512(tokenBytes)

	token := models.Token{
		ApplicationId: applicationId,
		Value:         hashed[:],
	}

	if s.storage.Create(&token) != nil {
		return ""
	}
	return fmt.Sprintf("VaulGuard-%d-%s", token.ID, base64.RawURLEncoding.EncodeToString(tokenBytes))
}

func (s service) Verify(token string) (models.Application, bool) {
	var id uint
	var value string

	_, err := fmt.Sscanf(token, "VaulGuard-%d-%s", &id, &value)

	if err != nil {
		return models.Application{}, false
	}

	t, err := s.storage.Get(id)

	if err != nil {
		return models.Application{}, false
	}

	decodedValue, err := base64.RawURLEncoding.DecodeString(value)
	hashedToken := blake3.Sum512(decodedValue)

	if err != nil {
		return models.Application{}, false
	}

	return t.Application, subtle.ConstantTimeCompare(hashedToken[:], t.Value) == 1
}