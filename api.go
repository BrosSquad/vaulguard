package main

import (
	"context"
	"github.com/BrosSquad/vaulguard/config"
	"github.com/BrosSquad/vaulguard/handlers"
	"github.com/BrosSquad/vaulguard/middleware"
	"github.com/BrosSquad/vaulguard/services"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	"log"
)

func registerAPIHandlers(ctx context.Context, cfg *config.Config, client *mongo.Database, db *gorm.DB, app *fiber.App) {
	var tokenCollection *mongo.Collection
	var secretCollection *mongo.Collection
	var applicationCollection *mongo.Collection
	apiV1 := app.Group("/api/v1")
	encryptionService, err := services.NewEncryptionService(cfg.ApplicationKey)

	if client != nil {
		tokenCollection = client.Collection("tokens")
		secretCollection = client.Collection("secrets")
		applicationCollection = client.Collection("applications")
	}

	tokenService := createTokenService(ctx, db, tokenCollection, cfg.StoreInSql)
	_ = createApplicationService(db, applicationCollection, cfg.StoreInSql)
	secretService := createSecretService(db, secretCollection, encryptionService, cfg.StoreInSql)

	apiV1.Use(middleware.TokenAuth(middleware.TokenAuthConfig{
		TokenService: tokenService,
		Header:       "authorization",
		HeaderPrefix: "token ",
	}))

	if err != nil {
		log.Fatalf("Cannot create encryption service: %v", err)
	}

	handlers.RegisterSecretHandlers(secretService, apiV1.Group("/secrets"))
}
