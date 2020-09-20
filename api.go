package main

import (
	"context"
	"github.com/BrosSquad/vaulguard/handlers"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	"log"

	"github.com/BrosSquad/vaulguard/config"
	"github.com/BrosSquad/vaulguard/middleware"
	"github.com/BrosSquad/vaulguard/services"
)

func registerAPIHandlers(ctx context.Context, cfg *config.Config, client *mongo.Database, db *gorm.DB, app *fiber.App) {
	var tokenCollection *mongo.Collection
	var secretCollection *mongo.Collection
	var applicationCollection *mongo.Collection
	apiV1 := app.Group("/api/v1")
	encryptionService, err := services.NewSecretKeyEncryption(cfg.ApplicationKey)

	if client != nil {
		tokenCollection = client.Collection("tokens")
		secretCollection = client.Collection("secrets")
		applicationCollection = client.Collection("applications")
	}

	tokenService := createTokenService(ctx, db, tokenCollection, cfg.UseSql)
	applicationService := createApplicationService(db, applicationCollection, cfg.UseSql)
	secretService := createSecretService(db, secretCollection, encryptionService, cfg.UseSql)
	if err != nil {
		log.Fatalf("Cannot create encryption service: %v", err)
	}

	secretsGroup := apiV1.Group("/secrets")

	secretsGroup.Use(middleware.TokenAuth(middleware.TokenAuthConfig{
		TokenService: tokenService,
		Header:       "authorization",
		HeaderPrefix: "token ",
	}))

	// TODO: Add JWT authentication

	handlers.RegisterSecretHandlers(secretService, secretsGroup)
	handlers.RegisterApplicationHandlers(applicationService, apiV1.Group("/applications"))
}
