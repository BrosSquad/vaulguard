package main

import (
	"context"
	"crypto/rand"
	"github.com/BrosSquad/vaulguard/config"
	"github.com/BrosSquad/vaulguard/services"
	"github.com/BrosSquad/vaulguard/services/application"
	"github.com/BrosSquad/vaulguard/services/secret"
	"github.com/BrosSquad/vaulguard/services/token"
	"github.com/gofiber/session/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"

	fsession "github.com/fasthttp/session/v2"
	"github.com/gofiber/session/v2/provider/redis"
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

func sessionIdGenerator() []byte {
	data := make([]byte, 32)
	n, err := rand.Read(data)

	if err != nil || n != 32 {
		return nil
	}

	return data
}

func createHttpSession(cfg *config.Config) *session.Session {
	var (
		provider fsession.Provider
		err      error
	)

	switch cfg.Http.Session.Provider {
	case "redis":
		provider, err = redis.New(redis.Config{
			KeyPrefix: "vaulguard_session",
			Addr:      cfg.Databases.Redis.Addr,
			Password:  cfg.Databases.Redis.Password,
			DB:        int(cfg.Http.Session.RedisDB),
		})

		if err != nil {
			return nil
		}
	default:
		provider = nil
	}

	return session.New(session.Config{
		Lookup:     "cookie:" + cfg.Http.Session.CookieName,
		Secure:     cfg.Http.Session.Secure,
		Domain:     cfg.Http.Session.Domain,
		SameSite:   cfg.Http.Session.SameSite,
		Expiration: cfg.Http.Session.Expiration,
		Provider:   provider,
		Generator:  sessionIdGenerator,
		GCInterval: cfg.Http.Session.GC,
	})
}
