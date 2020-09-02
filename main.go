package main

import (
	"context"
	"log"

	"github.com/BrosSquad/vaulguard/config"
	"github.com/BrosSquad/vaulguard/handlers"
	"github.com/BrosSquad/vaulguard/middleware"
	"github.com/BrosSquad/vaulguard/services"
	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

func registerAPIHandlers(cfg *config.Config, client *mongo.Client, db *gorm.DB, app *fiber.App) {
	apiV1 := app.Group("/api/v1")
	encryptionService, err := services.NewEncryptionService(cfg.ApplicationKey)
	tokenService := services.NewTokenService(db)

	apiV1.Use(middleware.TokenAuthMiddleware(middleware.TokenAuthConfig{
		TokenService: tokenService,
		Header:       "authorization",
		HeaderPrefix: "token ",
	}))

	if err != nil {
		log.Fatalf("Cannot create encryption service: %v", err)
	}

	secretService := createSecretService(db, client, encryptionService, cfg)
	handlers.RegisterSecretHandlers(secretService, apiV1.Group("/secrets"))
}

func main() {
	cfg, err := config.NewConfig()

	if err != nil {
		log.Fatalf("Error while creating app configuration: %v\n", err)
	}

	db, dbClose := connectToRelationalDatabaseAndMigrate(&cfg)
	defer dbClose()

	ctx, cancel := context.WithCancel(context.Background())
	client, mongoClose := connectToMongo(ctx, &cfg)
	defer mongoClose()

	app := fiber.New(&fiber.Settings{})

	registerAPIHandlers(&cfg, client, db, app)

	if err := app.Listen(cfg.Port); err != nil {
		log.Fatalf("Error while starting Fiber Server: %v", err)
	}

	cancel()
}
