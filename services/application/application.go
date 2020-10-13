package application

import (
	"context"
	"github.com/BrosSquad/vaulguard/models"
)

type Service interface {
	List(ctx context.Context, cb func([]models.ApplicationDto) error) error
	GetByName(ctx context.Context, name string) (models.ApplicationDto, error)
	Create(ctx context.Context, name string) (models.ApplicationDto, error)
	Get(ctx context.Context, page, perPage int) ([]models.ApplicationDto, error)
	GetOne(ctx context.Context, id interface{}) (models.ApplicationDto, error)
	Update(ctx context.Context, id interface{}, name string) (models.ApplicationDto, error)
	Delete(ctx context.Context, id interface{}) error
}
