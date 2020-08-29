package db

import (
	"context"
	"time"

	"github.com/BrosSquad/vaulguard/models"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DbConn      *gorm.DB
	MongoClient *mongo.Client
)

// ConnectToMongo - Connects to the running mongo database instance
func ConnectToMongo(ctx context.Context, url string) (_ *mongo.Client, err error) {
	MongoClient, err = mongo.NewClient(options.Client().ApplyURI(url))

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

func Migrate(useSqlSecretStorage bool) error {
	dst := []interface{}{
		&models.Application{},
		&models.Token{},
	}

	if useSqlSecretStorage {
		dst = append(dst, &models.Secret{})
	}

	return DbConn.AutoMigrate(dst...)
}

// ConnectToDatabaseProvider - Connects to different database providers supported by the application
// Supported providers:
// 1. PostgreSQL
// 2. MySQL
// 3. SQLite
func ConnectToDatabaseProvider(provider string, dsn string) (_ *gorm.DB, err error) {
	switch provider {
	case "postgres":
		return connectToPostgreSQL(dsn)
	case "mysql":
		return connectToMySQL(dsn)
	case "sqlite":
		return connectToSQLite(dsn)
	}

	return nil, nil
}

// ConnectToPostgres - Connects to the running postgres database instance
func connectToPostgreSQL(dsn string) (_ *gorm.DB, err error) {
	config := &gorm.Config{}
	DbConn, err = gorm.Open(postgres.New(postgres.Config{}), config)
	return DbConn, err
}

func connectToMySQL(dsn string) (_ *gorm.DB, err error) {
	return DbConn, nil
}

func connectToSQLite(dsn string) (_ *gorm.DB, err error) {
	DbConn, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	return DbConn, err
}
