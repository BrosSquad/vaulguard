package application

import "github.com/BrosSquad/vaulguard/models"

type Service interface {
	List(cb func([]models.Application) error) error
	GetByName(name string) (models.Application, error)
	Create(name string) (models.Application, error)
	Get(page, perPage int) ([]models.Application, error)
	GetOne(id uint) (models.Application, error)
	Update(id uint, name string) (models.Application, error)
	Delete(id uint) error
}
