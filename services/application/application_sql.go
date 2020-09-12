package application

import (
	"github.com/BrosSquad/vaulguard/models"
	"github.com/BrosSquad/vaulguard/services"
	"gorm.io/gorm"
)

type sqlService struct {
	db *gorm.DB
}

func NewSqlService(db *gorm.DB) Service {
	return sqlService{db: db}
}

func (a sqlService) List(cb func([]models.Application) error) error {
	results := make([]models.Application, 0, 50)
	err := a.db.FindInBatches(&results, 50, func(tx *gorm.DB, batch int) error {
		return cb(results)
	}).Error

	if err != nil {
		return err
	}

	return nil
}

func (a sqlService) GetByName(name string) (models.Application, error) {
	var app models.Application

	if err := a.db.Where("name = ?", name).Limit(1).Find(&app).Error; err != nil {
		return models.Application{}, err
	}

	return app, nil
}

func (a sqlService) Get(page, perPage int) ([]models.Application, error) {
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

func (a sqlService) GetOne(id uint) (models.Application, error) {
	app := models.Application{}

	tx := a.db.Find(&app, id)

	if tx.Error != nil {
		return app, tx.Error
	}

	return app, nil
}

func (a sqlService) Create(name string) (models.Application, error) {
	app := models.Application{}
	var count int64
	tx := a.db.Model(&app).Where("name = ?", name).Count(&count)

	if tx.Error != nil {
		return app, nil
	}

	if count > 0 {
		return app, services.ErrAlreadyExists
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

func (a sqlService) Update(id uint, name string) (models.Application, error) {
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

func (a sqlService) Delete(id uint) error {
	app := models.Application{}
	return a.db.Unscoped().Delete(&app, id).Error
}
