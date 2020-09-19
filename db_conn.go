package main

import (
	"context"
	"io"
	"log"

	"github.com/BrosSquad/vaulguard/db"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type mongoClose struct {
	ctx    context.Context
	client *mongo.Client
}

type gormClose struct {
	db *gorm.DB
}

func (g gormClose) Close() error {
	return db.Close()
}

func (m mongoClose) Close() error {
	return m.client.Disconnect(m.ctx)
}

func connectToMongoAndMigrate(ctx context.Context, mongoURI string) (*mongo.Client, io.Closer) {
	client, err := db.ConnectToMongo(ctx, mongoURI)

	if err != nil {
		log.Fatalf("Error while connecting to mongo db instance: %v\n", err)
	}

	if err := db.MongoCreateCollections(ctx, client); err != nil {
		log.Fatalf("Error while creating mongo collections: %v\n", err)
	}

	return client, &mongoClose{ctx, client}
}

func connectToRelationalDatabaseAndMigrate(database, dsn string) (*gorm.DB, io.Closer) {
	conn, err := db.ConnectToDatabaseProvider(database, dsn)

	if err != nil {
		log.Fatalf("Error while connection to PostgreSQL: %v\n", err)
	}

	if err := db.SqlMigrate(); err != nil {
		log.Fatalf("Auto migration failed: %v", err)
	}

	return conn, &gormClose{conn}
}
