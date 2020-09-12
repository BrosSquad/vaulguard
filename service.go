package main

import (
	"github.com/BrosSquad/vaulguard/config"
	"github.com/BrosSquad/vaulguard/services"
	"github.com/BrosSquad/vaulguard/services/application"
	"github.com/BrosSquad/vaulguard/services/secret"
	"github.com/BrosSquad/vaulguard/services/token"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

func createSecretService(db *gorm.DB, client *mongo.Client, encryption services.EncryptionService, cfg *config.Config) secret.Service {
	if cfg.StoreInSql {
		return secret.NewGormSecretStorage(db, encryption)
	}

	return secret.NewMongoClient(client, encryption)
}

func createApplicationService(db *gorm.DB, client *mongo.Client, cfg *config.Config) application.Service {
	if cfg.StoreInSql {
		return application.NewSqlService(db)
	}

	return application.NewMongoService(client)
}

func createTokenService(db *gorm.DB, client *mongo.Client, cfg *config.Config) token.Service {
	var storage token.Storage

	if cfg.StoreInSql {
		storage = token.NewSqlStorage(db)
	} else {
		storage = token.NewMongoStorage(client)
	}

	return token.NewService(storage)
}
