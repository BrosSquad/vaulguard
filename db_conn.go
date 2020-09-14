package main

import (
	"context"
	"io"
	"log"

	"github.com/BrosSquad/vaulguard/config"
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

func connectToMongo(ctx context.Context, cfg *config.Config) (*mongo.Client, io.Closer) {
	if cfg.StoreInSql {
		return nil, nil
	}
	client, err := db.ConnectToMongo(ctx, cfg.Mongo)

	if err != nil {
		log.Fatalf("Error while connecting to mongo db instance: %v\n", err)
	}

	return client, mongoClose{ctx, client}
}

func connectToRelationalDatabaseAndMigrate(cfg *config.Config) (*gorm.DB, io.Closer) {
	conn, err := db.ConnectToDatabaseProvider(cfg.Database, cfg.DatabaseDSN)

	if err != nil {
		log.Fatalf("Error while connection to PostgreSQL: %v", err)
	}

	if err := db.Migrate(); err != nil {
		log.Fatalf("Auto migration failed: %v", err)
	}

	return conn, gormClose{conn}
}
