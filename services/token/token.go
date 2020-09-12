package token

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/BrosSquad/vaulguard/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"lukechampine.com/blake3"
	"strconv"
	"strings"
	"time"
)

type Service interface {
	Generate(interface{}) string
	Verify(string) (models.ApplicationDto, bool)
}

type service struct {
	storage Storage
}

func NewService(storage Storage) Service {
	return service{storage}
}

func (s service) Generate(applicationId interface{}) string {
	tokenBytes := make([]byte, 64)
	_, err := rand.Read(tokenBytes)

	if err != nil {
		return ""
	}

	hashed := blake3.Sum512(tokenBytes)

	token, err := s.storage.Create(&models.TokenDto{
		ApplicationId: applicationId,
		Value:         hashed[:],
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	})

	if err != nil {
		return ""
	}

	var id string

	switch t := token.ID.(type) {
	case uint:
		id = strconv.FormatUint(uint64(t), 10)
	case primitive.ObjectID:
		id = t.Hex()
	}

	return fmt.Sprintf("VaulGuard.%s.%s", id, base64.RawURLEncoding.EncodeToString(tokenBytes))
}

func (s service) Verify(token string) (models.ApplicationDto, bool) {
	values := strings.Split(token, ".")

	if len(values) != 3 || values[0] != "VaulGuard" {
		return models.ApplicationDto{}, false
	}

	var id interface{}

	id, err := strconv.ParseUint(values[1], 10, 64)

	if err != nil {
		h, err := hex.DecodeString(values[1])

		if err != nil {
			return models.ApplicationDto{}, false
		}

		idObject := primitive.ObjectID{}
		if err := idObject.UnmarshalJSON(h); err != nil {
			return models.ApplicationDto{}, false
		}

		id = idObject
	}

	t, err := s.storage.Get(id)

	if err != nil {
		return models.ApplicationDto{}, false
	}

	decodedValue, err := base64.RawURLEncoding.DecodeString(values[2])
	hashedToken := blake3.Sum512(decodedValue)

	if err != nil {
		return models.ApplicationDto{}, false
	}

	return t.Application, subtle.ConstantTimeCompare(hashedToken[:], t.Value) == 1
}
