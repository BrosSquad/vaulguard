package secret

import (
	"context"
	"sync"

	"github.com/BrosSquad/vaulguard/models"
	"github.com/BrosSquad/vaulguard/services"
)

type Secret struct {
	Key   string
	Value string
}

type Service interface {
	Paginate(ctx context.Context, applicationID interface{}, page, perPage int) (map[string]string, error)
	Get(ctx context.Context, applicationID interface{}, key []string) (map[string]string, error)
	GetOne(ctx context.Context, applicationID interface{}, key string) (Secret, error)
	Create(ctx context.Context, applicationID interface{}, key, value string) (models.Secret, error)
	Update(ctx context.Context, applicationID interface{}, key, newKey, value string) (models.Secret, error)
	Delete(ctx context.Context, applicationID interface{}, key string) error
	InvalidateCache(ctx context.Context, applicationID interface{}) error
}

type baseService struct {
	mutex             *sync.RWMutex
	cacheLimit        int
	cache             [1024]map[string]models.Secret
	encryptionService services.Encryption
}
