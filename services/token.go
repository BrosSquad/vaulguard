package services

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/BrosSquad/vaulguard/models"
	"gorm.io/gorm"
	"lukechampine.com/blake3"
)

type TokenService interface {
	Generate(applicationId uint) string
	Verify() bool
}

type tokenService struct {
	db *gorm.DB
}

func NewTokenService(db *gorm.DB) TokenService {
	return tokenService{db: db}
}

func (s tokenService) Generate(applicationId uint) string {
	tokenBytes := make([]byte, 64)
	rand.Read(tokenBytes)
	token := models.Token{
		ApplicationId: applicationId,
		Value:         blake3.Sum512(tokenBytes),
	}
	s.db.Create(&token)

	return base64.RawURLEncoding.EncodeToString(tokenBytes)
}

func (s tokenService) Verify() bool {
	return false
}
