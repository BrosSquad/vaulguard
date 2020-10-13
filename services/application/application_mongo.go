package application

import (
	"context"
	"github.com/BrosSquad/vaulguard/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoService struct {
	client *mongo.Collection
}

func (m mongoService) List(ctx context.Context, cb func([]models.ApplicationDto) error) error {
	panic("implement me")
}

func (m mongoService) GetByName(ctx context.Context, name string) (models.ApplicationDto, error) {
	panic("implement me")
}

func (m mongoService) Create(ctx context.Context, name string) (models.ApplicationDto, error) {
	panic("implement me")
}

func (m mongoService) Get(ctx context.Context, page, perPage int) ([]models.ApplicationDto, error) {
	panic("implement me")
}

func (m mongoService) GetOne(ctx context.Context, id interface{}) (models.ApplicationDto, error) {
	panic("implement me")
}

func (m mongoService) Update(ctx context.Context, id interface{}, name string) (models.ApplicationDto, error) {
	panic("implement me")
}

func (m mongoService) Delete(ctx context.Context, id interface{}) error {
	panic("implement me")
}

func NewMongoService(client *mongo.Collection) Service {
	return mongoService{
		client: client,
	}
}
