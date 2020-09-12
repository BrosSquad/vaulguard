package main

import (
	"context"
	"log"

	"github.com/BrosSquad/vaulguard/config"
	"github.com/BrosSquad/vaulguard/db"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

func connectToMongo(ctx context.Context, cfg *config.Config) (*mongo.Client, func()) {
	if cfg.StoreInSql {
		return nil, func() {}
	}
	client, err := db.ConnectToMongo(ctx, cfg.Mongo)

	if err != nil {
		log.Fatalf("Error while connecting to mongo db instance: %v\n", err)
	}

	return client, func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatalf("Error while disconecting from mongodb instance: %v\n", err)
		}
	}
}

func connectToRelationalDatabaseAndMigrate(cfg *config.Config) (*gorm.DB, func() error) {
	conn, err := db.ConnectToDatabaseProvider(cfg.Database, cfg.DatabaseDSN)

	if err != nil {
		log.Fatalf("Error while connection to PostgreSQL: %v", err)
	}

	if err := db.Migrate(); err != nil {
		log.Fatalf("Auto migration failed: %v", err)
	}

	return conn, db.Close
}
