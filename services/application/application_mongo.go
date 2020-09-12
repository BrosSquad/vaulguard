package application

import (
	"github.com/BrosSquad/vaulguard/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoService struct {
	client *mongo.Client
}

func NewMongoService(client *mongo.Client) Service {
	return mongoService{
		client: client,
	}
}

func (m mongoService) List(cb func([]models.Application) error) error {
	panic("implement me")
}

func (m mongoService) GetByName(name string) (models.Application, error) {
	panic("implement me")
}

func (m mongoService) Create(name string) (models.Application, error) {
	panic("implement me")
}

func (m mongoService) Get(page, perPage int) ([]models.Application, error) {
	panic("implement me")
}

func (m mongoService) GetOne(id uint) (models.Application, error) {
	panic("implement me")
}

func (m mongoService) Update(id uint, name string) (models.Application, error) {
	panic("implement me")
}

func (m mongoService) Delete(id uint) error {
	panic("implement me")
}
