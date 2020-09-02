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
	Get(applicationID uint, page, perPage int) ([]Secret, error)
	GetOne(applicationID uint, key string) (Secret, error)
	Create(applicationID uint, key, value string) (models.Secret, error)
	Update(applicationID uint, key, newKey, value string) (models.Secret, error)
	Delete(applicationID uint, key string) error
}

type CacheKey struct {
	applicationID uint
	key           string
}

type gormSecretService struct {
	cacheLimit        int
	cache             map[CacheKey]models.Secret
	db                *gorm.DB
	encryptionService EncryptionService
}

func NewGormSecretStorage(db *gorm.DB, service EncryptionService) SecretService {
	return gormSecretService{
		cache:             make(map[CacheKey]models.Secret),
		cacheLimit:        1024,
		db:                db,
		encryptionService: service,
	}
}

func (g gormSecretService) Get(applicationID uint, page, perPage int) ([]Secret, error) {
	var secrets []models.Secret

	if page < 0 {
		page *= -1
	}

	err := g.db.
		Where("application_id = ?", applicationID).
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

func (g gormSecretService) GetOne(applicationID uint, key string) (Secret, error) {
	secret := models.Secret{}

	value, ok := g.cache[CacheKey{applicationID, key}]

	if ok {
		secret = value
	} else {
		err := g.db.Where("key = ? AND application_id = ?", key, applicationID).First(&secret).Error
		if err != nil {
			return Secret{}, err
		}
		if len(g.cache) < g.cacheLimit {
			g.cache[CacheKey{applicationID, key}] = secret
		}
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

func (g gormSecretService) Create(applicationID uint, key, value string) (models.Secret, error) {
	var count int64
	var secret models.Secret

	if err := g.db.Model(&secret).Where("key = ? AND application_id = ?", key, applicationID).Count(&count).Error; err != nil {
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
	secret.ApplicationId = applicationID

	if err := g.db.Create(&secret).Error; err != nil {
		return models.Secret{}, err
	}

	if len(g.cache) < g.cacheLimit {
		g.cache[CacheKey{applicationID, key}] = secret
	}

	return secret, nil
}

func (g gormSecretService) Update(applicationID uint, key, newKey, value string) (models.Secret, error) {
	secret := models.Secret{}

	if err := g.db.Where("key = ? AND application_id = ?", key, applicationID).Find(&secret).Error; err != nil {
		return secret, err
	}

	encrypted, err := g.encryptionService.EncryptString(value)

	if err != nil {
		return models.Secret{}, err
	}

	secret.Value = encrypted
	secret.Key = newKey

	if err := g.db.Save(&secret).Error; err != nil {
		return models.Secret{}, err
	}

	// Invalidate the old cache and add the new value to the cache
	cacheKey := CacheKey{applicationID, key}
	if _, ok := g.cache[cacheKey]; ok {
		delete(g.cache, cacheKey)
		g.cache[CacheKey{applicationID, newKey}] = secret
	}

	return secret, nil
}

func (g gormSecretService) Delete(applicationID uint, key string) error {
	secret := models.Secret{}

	cacheKey := CacheKey{applicationID, key}

	if _, ok := g.cache[cacheKey]; ok {
		delete(g.cache, cacheKey)
	}

	if err := g.db.Where("key = ? AND application_id = ?", key, applicationID).Delete(&secret).Error; err != nil {
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
