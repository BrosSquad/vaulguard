package main

import (
	"github.com/BrosSquad/vaulguard/config"
	"github.com/BrosSquad/vaulguard/services"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

func createSecretService(db *gorm.DB, client *mongo.Client, encryption services.EncryptionService, cfg *config.Config) services.SecretService {
	if cfg.StoreSecretInSql {
		return services.NewGormSecretStorage(db, encryption)
	}

	return nil
}
