package application

import (
	"context"
	"github.com/BrosSquad/vaulguard/models"
	"github.com/BrosSquad/vaulguard/services"
	"gorm.io/gorm"
)

type sqlService struct {
	db *gorm.DB
}

const size = 50

func (s sqlService) List(ctx context.Context, cb func([]models.ApplicationDto) error) error {
	results := make([]models.Application, 0, size)
	appsDto := make([]models.ApplicationDto, 0, size)
	err := s.db.WithContext(ctx).FindInBatches(&results, size, func(tx *gorm.DB, batch int) error {
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

func (s sqlService) GetByName(ctx context.Context, name string) (models.ApplicationDto, error) {
	app := models.Application{}

	if err := s.db.WithContext(ctx).Where("name = ?", name).Limit(1).Find(&app).Error; err != nil {
		return models.ApplicationDto{}, err
	}

	return models.ApplicationDto{
		ID:        app.ID,
		Name:      app.Name,
		CreatedAt: app.CreatedAt,
		UpdatedAt: app.UpdatedAt,
	}, nil
}

func (s sqlService) Create(ctx context.Context, name string) (models.ApplicationDto, error) {
	app := models.Application{}
	var count int64
	tx := s.db.WithContext(ctx).Model(&app).Where("name = ?", name).Count(&count)

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

func (s sqlService) Get(ctx context.Context, page, perPage int) ([]models.ApplicationDto, error) {
	if page < 0 {
		page *= -1
	}

	if perPage < 0 {
		perPage *= -1
	}

	apps := make([]models.Application, 0, perPage)

	if err := s.db.WithContext(ctx).Limit(perPage).Offset((page - 1) * perPage).Find(&apps).Error; err != nil {
		return nil, err
	}

	appsLen := len(apps)

	appsDto := make([]models.ApplicationDto, 0, appsLen)

	for _, app := range apps {
		appsDto = append(appsDto, models.ApplicationDto{
			ID:        app.ID,
			Name:      app.Name,
			CreatedAt: app.CreatedAt,
			UpdatedAt: app.UpdatedAt,
		})
	}

	return appsDto, nil
}

func (s sqlService) GetOne(ctx context.Context, id interface{}) (models.ApplicationDto, error) {
	app := models.Application{}

	if err := s.db.WithContext(ctx).First(&app, id.(uint)).Error; err != nil {
		return models.ApplicationDto{}, err
	}

	return models.ApplicationDto{
		ID:        app.ID,
		Name:      app.Name,
		CreatedAt: app.CreatedAt,
		UpdatedAt: app.UpdatedAt,
	}, nil
}

func (s sqlService) Update(ctx context.Context, id interface{}, name string) (models.ApplicationDto, error) {
	app := models.Application{}

	if err := s.db.WithContext(ctx).First(&app, id).Error; err != nil {
		return models.ApplicationDto{}, err
	}

	app.Name = name

	if err := s.db.WithContext(ctx).Save(&app).Error; err != nil {
		return models.ApplicationDto{}, err
	}

	return models.ApplicationDto{
		ID:        app.ID,
		Name:      app.Name,
		CreatedAt: app.CreatedAt,
		UpdatedAt: app.UpdatedAt,
	}, nil
}

func (s sqlService) Delete(ctx context.Context, id interface{}) error {
	app := models.Application{}
	return s.db.WithContext(ctx).Unscoped().Delete(&app, id).Error
}

func NewSqlService(db *gorm.DB) Service {
	return sqlService{db: db}
}
