package main

import (
	"context"
	"io"
	"log"

	"github.com/BrosSquad/vaulguard/config"
	"github.com/BrosSquad/vaulguard/handlers"
	"github.com/BrosSquad/vaulguard/middleware"
	"github.com/BrosSquad/vaulguard/services"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
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

	tokenService := createTokenService(ctx, db, tokenCollection, cfg)
	_ = createApplicationService(db, applicationCollection, cfg)
	secretService := createSecretService(db, secretCollection, encryptionService, cfg)

	apiV1.Use(middleware.TokenAuthMiddleware(middleware.TokenAuthConfig{
		TokenService: tokenService,
		Header:       "authorization",
		HeaderPrefix: "token ",
	}))

	if err != nil {
		log.Fatalf("Cannot create encryption service: %v", err)
	}

	handlers.RegisterSecretHandlers(secretService, apiV1.Group("/secrets"))
}

func main() {
	var db *gorm.DB
	var mongoClient *mongo.Client
	var mongoDatabase *mongo.Database

	var closer io.Closer
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Error while creating app configuration: %v\n", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	if cfg.StoreInSql {
		db, closer = connectToRelationalDatabaseAndMigrate(&cfg)
		defer closer.Close()
	} else {
		mongoClient, closer = connectToMongo(ctx, &cfg)
		mongoDatabase = mongoClient.Database("vaulguard")
		defer closer.Close()
	}
	app := fiber.New(fiber.Config{})

	registerAPIHandlers(ctx, &cfg, mongoDatabase, db, app)

	if err := app.Listen(cfg.Port); err != nil {
		log.Fatalf("Error while starting Fiber Server: %v", err)
	}

	cancel()
}
