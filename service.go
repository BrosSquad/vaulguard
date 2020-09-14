package main

import (
	"context"
	"github.com/BrosSquad/vaulguard/config"
	"github.com/BrosSquad/vaulguard/services"
	"github.com/BrosSquad/vaulguard/services/application"
	"github.com/BrosSquad/vaulguard/services/secret"
	"github.com/BrosSquad/vaulguard/services/token"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

func createSecretService(db *gorm.DB, client *mongo.Collection, encryption services.EncryptionService, cfg *config.Config) secret.Service {
	if cfg.StoreInSql {
		return secret.NewGormSecretStorage(secret.GormSecretConfig{
			Encryption: encryption,
			DB:         db,
		})
	}

	return secret.NewMongoClient(secret.MongoDBConfig{
		Encryption: encryption,
		Collection: client,
	})
}

func createApplicationService(db *gorm.DB, client *mongo.Collection, cfg *config.Config) application.Service {
	if cfg.StoreInSql {
		return application.NewSqlService(db)
	}
	return application.NewMongoService(client)
}

func createTokenService(ctx context.Context, db *gorm.DB, client *mongo.Collection, cfg *config.Config) token.Service {
	var storage token.Storage

	if cfg.StoreInSql {
		storage = token.NewSqlStorage(db)
	} else {
		storage = token.NewMongoStorage(ctx, client)
	}

	return token.NewService(storage)
}
