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
	Paginate(applicationID uint, page, perPage int) (map[string]string, error)
	Get(applicationID uint, key []string) (map[string]string, error)
	GetOne(applicationID uint, key string) (Secret, error)
	Create(applicationID uint, key, value string) (models.Secret, error)
	Update(applicationID uint, key, newKey, value string) (models.Secret, error)
	Delete(applicationID uint, key string) error
	InvalidateCache(applicationID uint) error
}

type baseService struct {
	mutex             *sync.RWMutex
	cacheLimit        int
	cache             map[uint]map[string]models.Secret
	encryptionService services.EncryptionService
}
