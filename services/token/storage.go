package token

import (
	"github.com/BrosSquad/vaulguard/models"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type Storage interface {
	Get(id uint) (models.Token, error)
	Create(token *models.Token) error
}

func NewSqlStorage(db *gorm.DB) Storage {
	return sqlStorage{db}
}

func NewMongoStorage(client *mongo.Client) Storage {
	return mongoStorage{client}
}

type sqlStorage struct {
	db *gorm.DB
}

func (s sqlStorage) Get(id uint) (models.Token, error) {
	token := models.Token{}

	tx := s.db.Joins("Application").First(&token, id)

	if err := tx.Error; err != nil {
		return models.Token{}, err
	}

	return token, nil
}

func (s sqlStorage) Create(token *models.Token) error {
	return s.db.Create(token).Error
}

type mongoStorage struct {
	client *mongo.Client
}

func (m mongoStorage) Get(id uint) (models.Token, error) {
	panic("implement me")
}

func (m mongoStorage) Create(token *models.Token) error {
	panic("implement me")
}
