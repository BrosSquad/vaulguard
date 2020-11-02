package main

import (
	"github.com/BrosSquad/vaulguard/services"
	"github.com/BrosSquad/vaulguard/services/application"
	"github.com/BrosSquad/vaulguard/services/secret"
	"github.com/BrosSquad/vaulguard/services/token"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

func createSecretService(db *gorm.DB, client *mongo.Collection, encryption services.Encryption, storeInSql bool) secret.Service {
	if storeInSql {
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

func createApplicationService(db *gorm.DB, client *mongo.Collection, storeInSql bool) application.Service {
	if storeInSql {
		return application.NewSqlService(db)
	}
	return application.NewMongoService(client)
}

func createTokenService(db *gorm.DB, client *mongo.Collection, storeInSql bool) token.Service {
	var storage token.Storage

	if storeInSql {
		storage = token.NewSqlStorage(db)
	} else {
		storage = token.NewMongoStorage(client)
	}

	return token.NewService(storage)
}
