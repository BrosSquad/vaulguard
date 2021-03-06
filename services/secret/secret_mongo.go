package secret

import (
	"context"
	"github.com/BrosSquad/vaulguard/models"
	"github.com/BrosSquad/vaulguard/services"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
)

type mongoService struct {
	baseService
	client *mongo.Collection
}

type MongoDBConfig struct {
	Encryption services.Encryption
	CacheSize  int
	Collection *mongo.Collection
}

func (m mongoService) Paginate(ctx context.Context, applicationID interface{}, page, perPage int) (map[string]string, error) {
	panic("implement me")
}

func (m mongoService) Get(ctx context.Context, applicationID interface{}, key []string) (map[string]string, error) {
	panic("implement me")
}

func (m mongoService) GetOne(ctx context.Context, applicationID interface{}, key string) (Secret, error) {
	panic("implement me")
}

func (m mongoService) Create(ctx context.Context, applicationID interface{}, key, value string) (models.Secret, error) {
	panic("implement me")
}

func (m mongoService) Update(ctx context.Context, applicationID interface{}, key, newKey, value string) (models.Secret, error) {
	panic("implement me")
}

func (m mongoService) Delete(ctx context.Context, applicationID interface{}, key string) error {
	panic("implement me")
}

func (m mongoService) InvalidateCache(ctx context.Context, applicationID interface{}) error {
	panic("implement me")
}

func NewMongoClient(config MongoDBConfig) Service {
	cacheSize := config.CacheSize

	if cacheSize == 0 {
		cacheSize = 8191
	}

	return &mongoService{
		baseService: baseService{
			mutex:             &sync.RWMutex{},
			cacheLimit:        cacheSize,
			cache:             [1024]map[string]models.Secret{},
			encryptionService: config.Encryption,
		},
		client: config.Collection,
	}
}
