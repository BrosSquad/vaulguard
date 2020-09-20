package secret

import (
	"log"
	"sync"

	"github.com/BrosSquad/vaulguard/models"
	"github.com/BrosSquad/vaulguard/services"
	"gorm.io/gorm"
)

type gormSecretService struct {
	baseService
	db *gorm.DB
}

type GormSecretConfig struct {
	Encryption services.Encryption
	CacheSize  int
	DB         *gorm.DB
}

func NewGormSecretStorage(config GormSecretConfig) Service {
	cacheSize := config.CacheSize

	if cacheSize == 0 {
		cacheSize = 8192
	}

	return &gormSecretService{
		baseService: baseService{
			mutex:             &sync.RWMutex{},
			cache:             [1024]map[string]models.Secret{},
			cacheLimit:        cacheSize,
			encryptionService: config.Encryption,
		},
		db: config.DB,
	}
}

func (g gormSecretService) Paginate(applicationID interface{}, page, perPage int) (map[string]string, error) {
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
		decryptedValue, err := g.encryptionService.DecryptString(s.Value)
		if err != nil {
			return nil, err
		}

		secretsDto[s.Key] = decryptedValue
	}

	return secretsDto, nil
}

func (g gormSecretService) GetOne(applicationID interface{}, key string) (Secret, error) {
	secret := models.Secret{}
	value, ok := g.cache[applicationID.(uint)][key]

	if ok {
		secret = value
	} else {
		err := g.db.Where("key = ? AND application_id = ?", key, applicationID).First(&secret).Error
		if err != nil {
			return Secret{}, err
		}
		go updateSecretCache(&g.baseService, []models.Secret{secret}, applicationID)
	}

	decryptedValue, err := g.encryptionService.DecryptString(secret.Value)

	if err != nil {
		return Secret{}, err
	}

	return Secret{
		Key:   key,
		Value: decryptedValue,
	}, nil
}

func updateSecretCache(g *baseService, secrets []models.Secret, applicationID interface{}) {

	if len(g.cache) >= g.cacheLimit {
		return
	}

	appId := applicationID.(uint)

	if m := g.cache[appId]; m == nil {
		g.cache[appId] = make(map[string]models.Secret, g.cacheLimit)
	}

	secretsMap := g.cache[appId]
	for _, s := range secrets {
		if _, ok := secretsMap[s.Key]; !ok && len(secretsMap) < g.cacheLimit {
			g.mutex.Lock()
			secretsMap[s.Key] = s
			g.mutex.Unlock()

		}
	}
}

func (g gormSecretService) Get(applicationID interface{}, keys []string) (_ map[string]string, err error) {
	var keysToFetch []string
	keysLen := len(keys)
	secrets := make([]models.Secret, 0, keysLen)

	for _, key := range keys {
		if s, ok := g.cache[applicationID.(uint)][key]; ok {
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

		go updateSecretCache(&g.baseService, secretsFetch, applicationID)
		for _, s := range secretsFetch {
			secrets = append(secrets, s)
		}
	}

	dtoSecrets := make(map[string]string, keysLen)

	for i := 0; i < len(secrets); i++ {
		decrypted, err := g.encryptionService.DecryptString(secrets[i].Value)
		if err != nil {
			return nil, err
		}

		dtoSecrets[keys[i]] = decrypted
	}

	return dtoSecrets, err
}

func (g gormSecretService) Create(applicationID interface{}, key, value string) (models.Secret, error) {
	var count int64
	var secret models.Secret

	if err := g.db.Model(&secret).Where("key = ? AND application_id = ?", key, applicationID).Count(&count).Error; err != nil {
		return secret, err
	}

	if count > 0 {
		return models.Secret{}, services.ErrAlreadyExists
	}

	encrypted, err := g.encryptionService.EncryptString(value)

	if err != nil {
		return models.Secret{}, err
	}

	secret.Key = key
	secret.Value = encrypted
	secret.ApplicationId = applicationID.(uint)

	if err := g.db.Create(&secret).Error; err != nil {
		return models.Secret{}, err
	}

	return secret, nil
}

func (g gormSecretService) Update(applicationID interface{}, key, newKey, value string) (models.Secret, error) {
	secret := models.Secret{}
	appId := applicationID.(uint)

	if err := g.db.Where("key = ? AND application_id = ?", key, appId).Find(&secret).Error; err != nil {
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
		if _, ok := g.cache[appId][key]; ok {
			g.mutex.Lock()
			delete(g.cache[appId], key)
			g.cache[appId][key] = secret
			g.mutex.Unlock()
		}
	}()

	return secret, nil
}

func (g gormSecretService) Delete(applicationID interface{}, key string) error {
	secret := models.Secret{}
	appId := applicationID.(uint)

	if _, ok := g.cache[appId][key]; ok {
		g.mutex.Lock()
		delete(g.cache[appId], key)
		g.mutex.Unlock()
	}

	if err := g.db.Where("key = ? AND application_id = ?", key, appId).Delete(&secret).Error; err != nil {
		return err
	}

	return nil
}

func (g gormSecretService) InvalidateCache(applicationID interface{}) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	appId := applicationID.(uint)
	g.cache[appId] = nil
	return nil
}
