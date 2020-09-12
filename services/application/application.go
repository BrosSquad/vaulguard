package application

import "github.com/BrosSquad/vaulguard/models"

type Service interface {
	List(cb func([]models.ApplicationDto) error) error
	GetByName(name string) (models.ApplicationDto, error)
	Create(name string) (models.ApplicationDto, error)
	Get(page, perPage int) ([]models.ApplicationDto, error)
	GetOne(id interface{}) (models.ApplicationDto, error)
	Update(id interface{}, name string) (models.ApplicationDto, error)
	Delete(id interface{}) error
}
