package services

import (
	"log"
	"sync"

	"github.com/BrosSquad/vaulguard/models"
	"gorm.io/gorm"
)

type Secret struct {
	Key   string
	Value string
}

type SecretService interface {
	Paginate(applicationID uint, page, perPage int) (map[string]string, error)
	Get(applicationID uint, key []string) (map[string]string, error)
	GetOne(applicationID uint, key string) (Secret, error)
	Create(applicationID uint, key, value string) (models.Secret, error)
	Update(applicationID uint, key, newKey, value string) (models.Secret, error)
	Delete(applicationID uint, key string) error
	InvalidateCache(applicationID uint) error
}
type secretService struct {
	mutex             *sync.RWMutex
	cacheLimit        int
	cache             map[uint]map[string]models.Secret
	encryptionService EncryptionService
}

type gormSecretService struct {
	secretService
	db *gorm.DB
}

func NewGormSecretStorage(db *gorm.DB, service EncryptionService) SecretService {
	return gormSecretService{
		secretService: secretService{
			mutex:             &sync.RWMutex{},
			cache:             make(map[uint]map[string]models.Secret, 1024),
			cacheLimit:        8192,
			encryptionService: service,
		},
		db: db,
	}
}

func (g gormSecretService) Paginate(applicationID uint, page, perPage int) (map[string]string, error) {
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

	secretsDto := make(map[string]string, len(secrets))

	for _, s := range secrets {
		decryptedValue, err := g.encryptionService.Decrypt(s.Value)
		if err != nil {
			return nil, err
		}

		secretsDto[s.Key] = decryptedValue
	}

	return secretsDto, nil
}

func (g gormSecretService) GetOne(applicationID uint, key string) (Secret, error) {
	secret := models.Secret{}

	value, ok := g.cache[applicationID][key]

	if ok {
		secret = value
	} else {
		err := g.db.Where("key = ? AND application_id = ?", key, applicationID).First(&secret).Error
		if err != nil {
			return Secret{}, err
		}
		go updateSecretCache(&g.secretService, []models.Secret{secret}, applicationID)
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

func updateSecretCache(g *secretService, secrets []models.Secret, applicationID uint) {

	if len(g.cache) >= g.cacheLimit {
		return
	}

	if _, ok := g.cache[applicationID]; !ok {
		g.cache[applicationID] = make(map[string]models.Secret, g.cacheLimit)
	}

	secretsMap := g.cache[applicationID]
	for _, s := range secrets {
		if _, ok := secretsMap[s.Key]; !ok && len(secretsMap) < g.cacheLimit {
			g.mutex.Lock()
			secretsMap[s.Key] = s
			g.mutex.Unlock()

		}
	}
}

func (g gormSecretService) Get(applicationID uint, keys []string) (_ map[string]string, err error) {
	var keysToFetch []string
	keysLen := len(keys)
	secrets := make([]models.Secret, 0, keysLen)

	for _, key := range keys {
		if s, ok := g.cache[applicationID][key]; ok {
			secrets = append(secrets, s)
		} else {
			keysToFetch = append(keysToFetch, key)
		}
	}

	if len(keysToFetch) > 0 {
		log.Printf("Keys to fetch: %d\n", len(keysToFetch))
		var secretsFetch []models.Secret
		result := g.db.Where("application_id = ? AND key IN ?", applicationID, keysToFetch).Find(&secretsFetch)

		if err = result.Error; err != nil {
			return nil, err
		}

		go updateSecretCache(&g.secretService, secretsFetch, applicationID)
		for _, s := range secretsFetch {
			secrets = append(secrets, s)
		}
	}

	dtoSecrets := make(map[string]string, keysLen)

	for i := 0; i < len(secrets); i++ {
		decrypted, err := g.encryptionService.Decrypt(secrets[i].Value)
		if err != nil {
			return nil, err
		}

		dtoSecrets[keys[i]] = decrypted
	}

	return dtoSecrets, err
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

	go func() {
		// Invalidate the old cache and add the new value to the cache
		if _, ok := g.cache[applicationID][key]; ok {
			g.mutex.Lock()
			delete(g.cache[applicationID], key)
			g.cache[applicationID][key] = secret
			g.mutex.Unlock()
		}
	}()

	return secret, nil
}

func (g gormSecretService) Delete(applicationID uint, key string) error {
	secret := models.Secret{}

	if _, ok := g.cache[applicationID][key]; ok {
		g.mutex.Lock()
		delete(g.cache[applicationID], key)
		g.mutex.Unlock()
	}

	if err := g.db.Where("key = ? AND application_id = ?", key, applicationID).Delete(&secret).Error; err != nil {
		return err
	}

	return nil
}

func (g gormSecretService) InvalidateCache(applicationID uint) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	for key := range g.cache[applicationID] {
		delete(g.cache[applicationID], key)
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
