package db

import (
	"context"
	"errors"
	"github.com/BrosSquad/vaulguard/log"
	"go.mongodb.org/mongo-driver/bson"
	"gorm.io/gorm/logger"
	"strings"
	"time"

	"github.com/BrosSquad/vaulguard/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	dbConn                          *gorm.DB
	ErrDatabaseProviderNotSupported = errors.New("database provider not supported (sqlite, postgres, mysql)")
)

type (
	Provider   int
	GormConfig struct {
		LogLevel    log.DBLevel
		Logger      logger.Writer
		SQLProvider Provider
		DSN         string
	}
)

const (
	MongoDBName                = "vaulguard"
	TokensMongoCollection      = "tokens"
	ApplicationMongoCollection = "applications"
	SecretsMongoCollection     = "secrets"

	PostgreSQL Provider = iota + 1
	MySQL
	SQLite
)

// ConnectToMongo - Connects to the running mongo database instance
func ConnectToMongo(ctx context.Context, url string) (_ *mongo.Client, err error) {
	MongoClient, err := mongo.NewClient(options.Client().ApplyURI(url))

	if err != nil {
		return nil, err
	}

	c, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err = MongoClient.Connect(c)

	if err != nil {
		return nil, err
	}

	return MongoClient, nil
}

func MongoCreateCollections(ctx context.Context, client *mongo.Client) error {
	database := client.Database(MongoDBName)

	if err := database.CreateCollection(ctx, TokensMongoCollection); err != nil {
		return err
	}

	if err := database.CreateCollection(ctx, ApplicationMongoCollection); err != nil {
		return err
	}

	applicationIndexes := []mongo.IndexModel{
		{
			Keys: bson.M{
				"Name": 1,
			},
			Options: options.Index().SetUnique(true),
		},
	}

	_, err := database.Collection(ApplicationMongoCollection).Indexes().CreateMany(ctx, applicationIndexes)

	if err != nil {
		return err
	}

	if err := database.CreateCollection(ctx, SecretsMongoCollection); err != nil {
		return err
	}

	secretIndexes := []mongo.IndexModel{
		{
			Keys: bson.M{
				"Key": 1,
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.M{
				"Value": 1,
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.M{
				"Key":           1,
				"ApplicationId": 1,
			},
		},
	}

	_, err = database.Collection(SecretsMongoCollection).Indexes().CreateMany(ctx, secretIndexes)

	if err != nil {
		return err
	}

	return nil
}

func SqlMigrate() error {
	dst := []interface{}{
		&models.Application{},
		&models.Token{},
		&models.Secret{},
	}

	return dbConn.AutoMigrate(dst...)
}

func GetDatabaseProvider(provider string) (Provider, error) {

	providerLower := strings.ToLower(provider)

	switch providerLower {
	case "mysql":
		return MySQL, nil
	case "postgres":
		return PostgreSQL, nil
	case "sqlite":
		return SQLite, nil

	}
	return 0, ErrDatabaseProviderNotSupported
}

// ConnectToDatabaseProvider - Connects to different database providers supported by the application
// Supported providers:
// 1. PostgreSQL
// 2. MySQL
// 3. SQLite
func ConnectToDatabaseProvider(config GormConfig) (_ *gorm.DB, err error) {
	gormConfig := &gorm.Config{
		PrepareStmt: true,
	}

	if config.Logger != nil {
		gormConfig.Logger = logger.New(config.Logger, logger.Config{
			Colorful: true,
			LogLevel: logger.LogLevel(config.LogLevel),
		})
	}

	switch config.SQLProvider {
	case PostgreSQL:
		return connectToPostgreSQL(config.DSN, gormConfig)
	case MySQL:
		return connectToMySQL(config.DSN, gormConfig)
	case SQLite:
		return connectToSQLite(config.DSN, gormConfig)
	}

	return nil, nil
}

// ConnectToPostgres - Connects to the running postgres database instance
func connectToPostgreSQL(dsn string, config *gorm.Config) (_ *gorm.DB, err error) {
	dbConn, err = gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}), config)
	return dbConn, err
}

func connectToMySQL(dsn string, config *gorm.Config) (_ *gorm.DB, err error) {
	return dbConn, nil
}

func connectToSQLite(dsn string, config *gorm.Config) (_ *gorm.DB, err error) {
	dbConn, err = gorm.Open(sqlite.Open(dsn), config)
	return dbConn, err
}

func Close() error {
	db, err := dbConn.DB()

	if err != nil {
		return err
	}

	return db.Close()
}
