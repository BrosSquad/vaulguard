package application

import (
	"github.com/BrosSquad/vaulguard/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoService struct {
	client *mongo.Collection
}

func (m mongoService) List(cb func([]models.ApplicationDto) error) error {
	panic("implement me")
}

func (m mongoService) GetByName(name string) (models.ApplicationDto, error) {
	panic("implement me")
}

func (m mongoService) Create(name string) (models.ApplicationDto, error) {
	panic("implement me")
}

func (m mongoService) Get(page, perPage int) ([]models.ApplicationDto, error) {
	panic("implement me")
}

func (m mongoService) GetOne(id interface{}) (models.ApplicationDto, error) {
	panic("implement me")
}

func (m mongoService) Update(id interface{}, name string) (models.ApplicationDto, error) {
	panic("implement me")
}

func (m mongoService) Delete(id interface{}) error {
	panic("implement me")
}

func NewMongoService(client *mongo.Collection) Service {
	return mongoService{
		client: client,
	}
}
