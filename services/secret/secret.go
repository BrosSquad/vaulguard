package secret

import (
	"sync"

	"github.com/BrosSquad/vaulguard/models"
	"github.com/BrosSquad/vaulguard/services"
)

type Secret struct {
	Key   string
	Value string
}

type Service interface {
	Paginate(applicationID interface{}, page, perPage int) (map[string]string, error)
	Get(applicationID interface{}, key []string) (map[string]string, error)
	GetOne(applicationID interface{}, key string) (Secret, error)
	Create(applicationID interface{}, key, value string) (models.Secret, error)
	Update(applicationID interface{}, key, newKey, value string) (models.Secret, error)
	Delete(applicationID interface{}, key string) error
	InvalidateCache(applicationID interface{}) error
}

type baseService struct {
	mutex             *sync.RWMutex
	cacheLimit        int
	cache             [1024]map[string]models.Secret
	encryptionService services.Encryption
}
