package services

import (
	"errors"

	"github.com/BrosSquad/vaulguard/models"
	"gorm.io/gorm"
)

var ErrAlreadyExists = errors.New("Model already exists.")

type ApplicationService interface {
	List(cb func([]models.Application) error) error
	GetByName(name string) (models.Application, error)
	Create(name string) (models.Application, error)
	Get(page, perPage int) ([]models.Application, error)
	GetOne(id uint) (models.Application, error)
	Update(id uint, name string) (models.Application, error)
	Delete(id uint) error
}

type applicationService struct {
	db *gorm.DB
}

func NewApplicationService(db *gorm.DB) ApplicationService {
	return applicationService{db: db}
}

func (a applicationService) List(cb func([]models.Application) error) error {
	results := make([]models.Application, 0, 50)
	err := a.db.FindInBatches(&results, 50, func(tx *gorm.DB, batch int) error {
		return cb(results)
	}).Error

	if err != nil {
		return err
	}

	return nil
}

func (a applicationService) GetByName(name string) (models.Application, error) {
	var app models.Application

	if err := a.db.Where("name = ?", name).Limit(1).Find(&app).Error; err != nil {
		return models.Application{}, err
	}

	return app, nil
}

func (a applicationService) Get(page, perPage int) ([]models.Application, error) {
	var apps []models.Application

	if page < 0 {
		page *= -1
	}

	tx := a.db.Limit(perPage).Offset((page - 1) * perPage).Find(&apps)

	if tx.Error != nil {
		return apps, tx.Error
	}

	return apps, nil
}

func (a applicationService) GetOne(id uint) (models.Application, error) {
	app := models.Application{}

	tx := a.db.Find(&app, id)

	if tx.Error != nil {
		return app, tx.Error
	}

	return app, nil
}

func (a applicationService) Create(name string) (models.Application, error) {
	app := models.Application{}
	var count int64
	tx := a.db.Model(&app).Where("name = ?", name).Count(&count)

	if tx.Error != nil {
		return app, nil
	}

	if count > 0 {
		return app, ErrAlreadyExists
	}

	application := models.Application{
		Name: name,
	}

	tx = a.db.Create(&application)

	if tx.Error != nil {
		return models.Application{}, tx.Error
	}

	return application, nil
}

func (a applicationService) Update(id uint, name string) (models.Application, error) {
	app := models.Application{}

	a.db.Model(&app).Where("id = ?", id).Update("name", name)

	tx := a.db.First(&app, id)

	if tx.Error != nil {
		return app, tx.Error
	}

	app.Name = name

	tx = a.db.Save(&app)

	if tx.Error != nil {
		return models.Application{}, tx.Error
	}

	return app, nil
}

func (a applicationService) Delete(id uint) error {
	app := models.Application{}
	return a.db.Unscoped().Delete(&app, id).Error
}
