package main

import (
	"context"
	"flag"
	"github.com/BrosSquad/vaulguard/config"
	"github.com/BrosSquad/vaulguard/db"
	"github.com/BrosSquad/vaulguard/handlers"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	"io"
	"log"
	"os"
)

func createConfig(configPath string) (*config.Config, error) {
	var err error

	configPath, err = getAbsolutePath(configPath)

	if err != nil {
		return nil, err
	}

	cfgFile, err := os.OpenFile(configPath, os.O_RDONLY, DefaultPermission)

	if err != nil {
		return nil, err
	}

	cfg, err := config.NewConfig(cfgFile)

	if err != nil {
		return nil, err
	}

	if err := cfgFile.Close(); err != nil {
		return nil, err
	}

	return cfg, err
}

func main() {
	var sqlDb *gorm.DB
	var mongoClient *mongo.Client
	var mongoDatabase *mongo.Database
	var closer io.Closer

	configPath := flag.String("config", "./config.yml", "Path to config file")
	flag.Parse()

	cfg, err := createConfig(*configPath)

	if err != nil {
		log.Fatalf("Error while creating application configuration: %v\n", err)
	}

	key, err := getKeys(cfg)

	if err != nil {
		log.Fatalf("Error while loading application keys: %v", err)
	}

	cfg.ApplicationKey = key

	ctx, cancel := context.WithCancel(context.Background())

	if cfg.UseSql {
		sqlDb, closer = connectToRelationalDatabaseAndMigrate(cfg.Databases.SQL.Provider, cfg.Databases.SQL.DSN)
		defer closer.Close()
	} else {
		mongoClient, closer = connectToMongoAndMigrate(ctx, cfg.Databases.Mongo.URI)
		mongoDatabase = mongoClient.Database(db.MongoDBName)
		defer closer.Close()
	}

	app := fiber.New(fiber.Config{
		Prefork:      cfg.Http.Prefork,
		ErrorHandler: handlers.Error,
	})

	registerAPIHandlers(ctx, cfg, mongoDatabase, sqlDb, app)

	if err := app.Listen(cfg.Http.Address); err != nil {
		log.Fatalf("Error while starting http server: %v", err)
	}

	cancel()
}
