package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
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

// ConnectToPostgres - Connects to the running postgres database instance
func ConnectToPostgres(dsn string) (_ *gorm.DB, err error) {
	config := &gorm.Config{}
	DbConn, err = gorm.Open(postgres.New(postgres.Config{}), config)
	return DbConn, err
}
