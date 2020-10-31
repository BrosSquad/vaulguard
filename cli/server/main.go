package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/session/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	"io"
	"log"
	"os"
	"os/signal"

	"github.com/BrosSquad/vaulguard/api"
	"github.com/BrosSquad/vaulguard/config"
	"github.com/BrosSquad/vaulguard/db"
	"github.com/BrosSquad/vaulguard/handlers"
	vaulguardlog "github.com/BrosSquad/vaulguard/log"
	"github.com/BrosSquad/vaulguard/services"
	"github.com/BrosSquad/vaulguard/utils"
)

func createConfig(configPath string, port int) (cfg *config.Config, err error) {
	configPath, err = utils.GetAbsolutePath(configPath)

	if err != nil {
		return nil, err
	}

	cfgFile, err := os.OpenFile(configPath, os.O_RDONLY, DefaultKeysPermission)

	if err != nil {
		return nil, err
	}

	cfg, err = config.New(cfgFile)

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
	var (
		sqlDb                 *gorm.DB
		mongoClient           *mongo.Client
		mongoDatabase         *mongo.Database
		tokenCollection       *mongo.Collection
		secretCollection      *mongo.Collection
		applicationCollection *mongo.Collection
		httpSession           *session.Session
		closer                io.Closer
	)

	signalCh := make(chan os.Signal, 1)
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

	ctx, cancel := context.WithCancel(context.Background())
	signal.Notify(signalCh, os.Interrupt)
	go func(cancel context.CancelFunc) {
		select {
		case <-signalCh:
			cancel()
		}
	}(cancel)

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
	english := en.New()
	uni := ut.New(english, english)
	englishTranslations, found := uni.GetTranslator(cfg.Locale)

	if !found {
		logger.Fatalf(errors.New("locale not found"), "No translations found for locale %s", cfg.Locale)
	}

	app := fiber.New(fiber.Config{
		Prefork:      cfg.Http.Prefork,
		ErrorHandler: handlers.Error(englishTranslations),
	})

	if cfg.Debug {
		logger.Debug("Adding pprof routes\n")
		app.Use(pprof.New())
		if cfg.MemoryUsage.Report {
			go utils.MemoryUsage(ctx, cfg.MemoryUsage.Sleep, logger)
		}
	}

	if cfg.UseDashboard {
		httpSession = createHttpSession(cfg)
	}

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
		Validator:             v,
		Session:               httpSession,
	}

	fiberAPI.RegisterHandlers()

	go func() {
		logger.Debug("Start to listen on: %s\n", cfg.Http.Address)
		logger.Debug("Preforking? %v", cfg.Http.Prefork)
		if err := app.Listen(cfg.Http.Address); err != nil {
			logger.Fatalf(err, "Error while starting http server\n")
		}
	}()

	<-ctx.Done()
	logger.Debug("Shutting down application...\n")
	if err := app.Shutdown(); err != nil {
		logger.Fatalf(err, "Error while shutting down the api\n")
	}
	logger.Debug("Exiting...\n")

}
