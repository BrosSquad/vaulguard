package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"io"
	"log"
	"os"

	"github.com/BrosSquad/vaulguard/api"
	"github.com/BrosSquad/vaulguard/config"
	"github.com/BrosSquad/vaulguard/db"
	"github.com/BrosSquad/vaulguard/handlers"
	vaulguardlog "github.com/BrosSquad/vaulguard/log"
	"github.com/BrosSquad/vaulguard/services"
	"github.com/BrosSquad/vaulguard/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

func createConfig(configPath string, port int) (*config.Config, error) {
	var err error

	configPath, err = utils.GetAbsolutePath(configPath)

	if err != nil {
		return nil, err
	}

	cfgFile, err := os.OpenFile(configPath, os.O_RDONLY, DefaultKeysPermission)

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

	if port != 0 {
		cfg.Http.Address = fmt.Sprintf(":%d", port)
	}

	return cfg, err
}

func main() {
	var sqlDb *gorm.DB
	var mongoClient *mongo.Client
	var mongoDatabase *mongo.Database
	var tokenCollection *mongo.Collection
	var secretCollection *mongo.Collection
	var applicationCollection *mongo.Collection
	var closer io.Closer

	configPath := flag.String("config", "./config.yml", "Path to config file")
	port := flag.Int("port", 0, "Default port, overrides usage from config")

	flag.Parse()

	cfg, err := createConfig(*configPath, *port)

	if err != nil {
		log.Fatalf("Error while creating application configuration: %v\n", err)
	}

	logger := vaulguardlog.NewVaulGuardLogger(vaulguardlog.GetLogLevel(cfg.Logging.Level), cfg.UseConsole)
	vaulguardlog.SetDefaultLogger(logger)

	key, err := getKeys(cfg)

	if err != nil {
		logger.Fatalf(err, "Error while loading application keys\n")
	}

	cfg.ApplicationKey = key

	// TODO: Handle os Signals
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
		tokenCollection = mongoDatabase.Collection("tokens")
		secretCollection = mongoDatabase.Collection("secrets")
		applicationCollection = mongoDatabase.Collection("applications")
		defer closer.Close()
	}

	encryptionService, err := services.NewSecretKeyEncryption(cfg.ApplicationKey)

	if err != nil {
		logger.Fatalf(err, "Error while creating encryption service\n")
	}

	v := validator.New()
	en := en.New
	translator := ut.New()

	app := fiber.New(fiber.Config{
		Prefork:      cfg.Http.Prefork,
		ErrorHandler: handlers.Error(),
	})

	fiberAPI := api.Fiber{
		Ctx:                   ctx,
		Cfg:                   cfg,
		App:                   app.Group("/api/v1"),
		TokenCollection:       tokenCollection,
		SecretCollection:      secretCollection,
		ApplicationCollection: applicationCollection,
		SecretService:         createSecretService(sqlDb, secretCollection, encryptionService, cfg.UseSql),
		ApplicationService:    createApplicationService(sqlDb, applicationCollection, cfg.UseSql),
		TokenService:          createTokenService(ctx, sqlDb, tokenCollection, cfg.UseSql),
		Logger:                logger,
	}

	fiberAPI.RegisterHandlers()

	go utils.MemoryUsage(ctx, logger)

	if err := app.Listen(cfg.Http.Address); err != nil {
		logger.Fatalf(err, "Error while starting http server\n")
	}

	cancel()
}
