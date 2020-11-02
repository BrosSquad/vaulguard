package application

import (
	"context"
	"github.com/BrosSquad/vaulguard/models"
)

type Service interface {
	List(context.Context, int, func([]models.ApplicationDto) error) error
	GetByName(context.Context, string) (models.ApplicationDto, error)
	Create(context.Context, string) (models.ApplicationDto, error)
	Get(context.Context, int, int) ([]models.ApplicationDto, error)
	GetOne(context.Context, interface{}) (models.ApplicationDto, error)
	Update(context.Context, interface{}, string) (models.ApplicationDto, error)
	Delete(context.Context, interface{}) error
}
