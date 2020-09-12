package secret

import (
	"github.com/BrosSquad/vaulguard/models"
	"github.com/BrosSquad/vaulguard/services"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
)

type mongoService struct {
	baseService
	client *mongo.Client
}


func NewMongoClient(client *mongo.Client, encryption services.EncryptionService) Service {
	return mongoService{
		baseService: baseService{
			mutex:             &sync.RWMutex{},
			cacheLimit:        8192,
			cache:             make(map[uint]map[string]models.Secret),
			encryptionService: encryption,
		},
		client:      client,
	}
}

func (m mongoService) Paginate(applicationID uint, page, perPage int) (map[string]string, error) {
	panic("implement me")
}

func (m mongoService) Get(applicationID uint, key []string) (map[string]string, error) {
	panic("implement me")
}

func (m mongoService) GetOne(applicationID uint, key string) (Secret, error) {
	panic("implement me")
}

func (m mongoService) Create(applicationID uint, key, value string) (models.Secret, error) {
	panic("implement me")
}

func (m mongoService) Update(applicationID uint, key, newKey, value string) (models.Secret, error) {
	panic("implement me")
}

func (m mongoService) Delete(applicationID uint, key string) error {
	panic("implement me")
}

func (m mongoService) InvalidateCache(applicationID uint) error {
	panic("implement me")
}
