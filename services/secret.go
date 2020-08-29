package services

import (
	"github.com/BrosSquad/vaulguard/models"
	"gorm.io/gorm"
)

type Secret struct {
	Key   string
	Value string
}

type SecretService interface {
	Get(applicationId uint, page, perPage int) ([]Secret, error)
	GetOne(applicationId uint, key string) (Secret, error)
	Create(applicationId uint, key, value string) (models.Secret, error)
	Update(applicationId uint, key, value string) (models.Secret, error)
	Delete(applicationId uint, key string) error
}

type gormSecretService struct {
	db                *gorm.DB
	encryptionService EncryptionService
}

func NewGormSecretStorage(db *gorm.DB, service EncryptionService) SecretService {
	return gormSecretService{
		db:                db,
		encryptionService: service,
	}
}

func (g gormSecretService) Get(applicationId uint, page, perPage int) ([]Secret, error) {
	var secrets []models.Secret

	if page < 0 {
		page *= -1
	}

	err := g.db.
		Where("application_id = ?", applicationId).
		Limit(perPage).
		Offset((page - 1) * perPage).
		Find(&secrets).Error

	if err != nil {
		return nil, err
	}

	secretsDto := make([]Secret, len(secrets))

	for i, s := range secrets {
		decryptedValue, err := g.encryptionService.Decrypt(s.Value)
		if err != nil {
			return nil, err
		}

		secretsDto[i] = Secret{
			Key:   s.Key,
			Value: decryptedValue,
		}
	}

	return secretsDto, nil
}

func (g gormSecretService) GetOne(applicationId uint, key string) (Secret, error) {
	secret := models.Secret{}

	if err := g.db.Where("key = ? AND application_id = ?", key, applicationId).First(&secret).Error; err != nil {
		return Secret{}, err
	}

	decryptedValue, err := g.encryptionService.Decrypt(secret.Value)

	if err != nil {
		return Secret{}, err
	}

	return Secret{
		Key:   key,
		Value: decryptedValue,
	}, nil
}

func (g gormSecretService) Create(applicationId uint, key, value string) (models.Secret, error) {
	var count int64
	var secret models.Secret

	if err := g.db.Model(&secret).Where("key = ? AND application_id = ?", key, applicationId).Count(&count).Error; err != nil {
		return secret, err
	}

	if count > 0 {
		return models.Secret{}, ErrAlreadyExists
	}

	encrypted, err := g.encryptionService.EncryptString(value)

	if err != nil {
		return models.Secret{}, err
	}

	secret.Key = key
	secret.Value = encrypted
	secret.ApplicationId = applicationId

	if err := g.db.Create(&secret).Error; err != nil {
		return models.Secret{}, err
	}

	return secret, nil
}

func (g gormSecretService) Update(applicationId uint, key, value string) (models.Secret, error) {
	secret := models.Secret{}

	if err := g.db.Where("key = ? AND application_id = ?", key, applicationId).Find(&secret).Error; err != nil {
		return secret, err
	}

	encrypted, err := g.encryptionService.EncryptString(value)

	if err != nil {
		return models.Secret{}, err
	}

	secret.Value = encrypted

	if err := g.db.Save(&secret).Error; err != nil {
		return models.Secret{}, err
	}

	return secret, nil
}

func (g gormSecretService) Delete(applicationId uint, key string) error {
	secret := models.Secret{}

	if err := g.db.Where("key = ? AND application_id = ?", key, applicationId).Delete(&secret).Error; err != nil {
		return err
	}

	return nil
}

//type mongoSecretService struct {
//	client            *mongo.Client
//	encryptionService EncryptionService
//}

//func NewMongoSecretStorage(client *mongo.Client, service EncryptionService) SecretService {
//	return mongoSecretService{
//		client:            client,
//		encryptionService: service,
//	}
//}
