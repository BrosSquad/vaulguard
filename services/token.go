package services

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"github.com/BrosSquad/vaulguard/models"
	"gorm.io/gorm"
	"lukechampine.com/blake3"
)

type TokenService interface {
	Generate(uint) string
	Verify(string) bool
}

type tokenService struct {
	db *gorm.DB
}

func NewTokenService(db *gorm.DB) TokenService {
	return tokenService{db: db}
}

func (s tokenService) Generate(applicationId uint) string {
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
	tx := s.db.Create(&token)

	if tx.Error != nil {
		return ""
	}
	return fmt.Sprintf("VaulGuard-%d-%s", token.ID, base64.RawURLEncoding.EncodeToString(tokenBytes))
}

func (s tokenService) Verify(token string) bool {
	var id uint
	var value string
	tokenModel := models.Token{}

	_, err := fmt.Sscanf(token, "VaulGuard-%d-%s", &id, &value)

	if err != nil {
		return false
	}

	tx := s.db.First(&tokenModel, id)

	if tx.Error != nil {
		return false
	}


	decodedValue, err := base64.RawURLEncoding.DecodeString(value)
	hashedToken := blake3.Sum512(decodedValue)

	if err != nil {
		return false
	}

	return subtle.ConstantTimeCompare(hashedToken[:], tokenModel.Value) == 1
}
