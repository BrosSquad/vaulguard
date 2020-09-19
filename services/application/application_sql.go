package application

import (
	"github.com/BrosSquad/vaulguard/models"
	"github.com/BrosSquad/vaulguard/services"
	"gorm.io/gorm"
)

type sqlService struct {
	db *gorm.DB
}

const size = 50

func (s sqlService) List(cb func([]models.ApplicationDto) error) error {
	results := make([]models.Application, 0, size)
	appsDto := make([]models.ApplicationDto, 0, size)
	err := s.db.FindInBatches(&results, size, func(tx *gorm.DB, batch int) error {
		appsDto = appsDto[:0]
		for _, result := range results {
			appsDto = append(appsDto, models.ApplicationDto{
				ID:        result.ID,
				Name:      result.Name,
				CreatedAt: result.CreatedAt,
				UpdatedAt: result.UpdatedAt,
			})
		}
		results = results[:0]
		return cb(appsDto)
	}).Error

	if err != nil {
		return err
	}

	return nil
}

func (s sqlService) GetByName(name string) (models.ApplicationDto, error) {
	app := models.Application{}

	if err := s.db.Where("name = ?", name).Limit(1).Find(&app).Error; err != nil {
		return models.ApplicationDto{}, err
	}

	return models.ApplicationDto{
		ID:        app.ID,
		Name:      app.Name,
		CreatedAt: app.CreatedAt,
		UpdatedAt: app.UpdatedAt,
	}, nil
}

func (s sqlService) Create(name string) (models.ApplicationDto, error) {
	app := models.Application{}
	var count int64
	tx := s.db.Model(&app).Where("name = ?", name).Count(&count)

	if tx.Error != nil {
		return models.ApplicationDto{}, nil
	}

	if count > 0 {
		return models.ApplicationDto{}, services.ErrAlreadyExists
	}

	app = models.Application{
		Name: name,
	}

	tx = s.db.Create(&app)

	if tx.Error != nil {
		return models.ApplicationDto{}, tx.Error
	}

	return models.ApplicationDto{
		ID:        app.ID,
		Name:      app.Name,
		CreatedAt: app.CreatedAt,
		UpdatedAt: app.UpdatedAt,
	}, nil
}

func (s sqlService) Get(page, perPage int) ([]models.ApplicationDto, error) {
	apps := make([]models.Application, 0, perPage)
	if page < 0 {
		page *= -1
	}

	tx := s.db.Limit(perPage).Offset((page - 1) * perPage).Find(&apps)

	if tx.Error != nil {
		return nil, tx.Error
	}

	appsLen := len(apps)

	appsDto := make([]models.ApplicationDto, appsLen)

	for i := 0; i < appsLen; i++ {
		appsDto[i] = models.ApplicationDto{
			ID:        apps[i].ID,
			Name:      apps[i].Name,
			CreatedAt: apps[i].CreatedAt,
			UpdatedAt: apps[i].UpdatedAt,
		}
	}

	return appsDto, nil
}

func (s sqlService) GetOne(id interface{}) (models.ApplicationDto, error) {
	app := models.Application{}

	tx := s.db.Find(&app, id.(uint))

	if tx.Error != nil {
		return models.ApplicationDto{}, tx.Error
	}

	return models.ApplicationDto{
		ID:        app.ID,
		Name:      app.Name,
		CreatedAt: app.CreatedAt,
		UpdatedAt: app.UpdatedAt,
	}, nil
}

func (s sqlService) Update(id interface{}, name string) (models.ApplicationDto, error) {
	app := models.Application{}

	tx := s.db.First(&app, id)

	if tx.Error != nil {
		return models.ApplicationDto{}, tx.Error
	}

	app.Name = name

	tx = s.db.Save(&app)

	if tx.Error != nil {
		return models.ApplicationDto{}, tx.Error
	}

	return models.ApplicationDto{
		ID:        app.ID,
		Name:      app.Name,
		CreatedAt: app.CreatedAt,
		UpdatedAt: app.UpdatedAt,
	}, nil
}

func (s sqlService) Delete(id interface{}) error {
	app := models.Application{}
	return s.db.Unscoped().Delete(&app, id).Error
}

func NewSqlService(db *gorm.DB) Service {
	return sqlService{db: db}
}