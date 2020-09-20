package main

import (
	"context"
	"flag"
	"github.com/BrosSquad/vaulguard/config"
	"github.com/BrosSquad/vaulguard/db"
	"github.com/BrosSquad/vaulguard/handlers"
	vaulguardlog "github.com/BrosSquad/vaulguard/log"
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

	logger := vaulguardlog.NewVaulGuardLogger(vaulguardlog.GetLogLevel(cfg.Logging.Level), cfg.UseConsole)

	key, err := getKeys(cfg)

	if err != nil {
		logger.Fatalf(err, "Error while loading application keys\n")
	}

	cfg.ApplicationKey = key

	ctx, cancel := context.WithCancel(context.Background())

	if cfg.UseSql {
		sqlDb, closer, err = connectToRelationalDatabaseAndMigrate(logger, cfg)
		if err != nil {
			logger.Fatalf(err, "Error while connecting to database\n")
		}
		defer closer.Close()
	} else {
		mongoClient, closer, err = connectToMongoAndMigrate(ctx, cfg.Databases.Mongo.URI)
		if err != nil {
			logger.Fatalf(err, "Error while connecting to MongoDB\n")
		}
		mongoDatabase = mongoClient.Database(db.MongoDBName)
		defer closer.Close()
	}

	app := fiber.New(fiber.Config{
		Prefork:      cfg.Http.Prefork,
		ErrorHandler: handlers.Error,
	})

	registerAPIHandlers(ctx, cfg, mongoDatabase, sqlDb, app)

	if err := app.Listen(cfg.Http.Address); err != nil {
		logger.Fatalf(err, "Error while starting http server\n")
	}

	cancel()
}
