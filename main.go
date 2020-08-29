package main

import (
	"context"
	"log"

	"github.com/BrosSquad/vaulguard/config"
	"github.com/BrosSquad/vaulguard/db"
	"go.mongodb.org/mongo-driver/mongo"
)

func connectToMongo(ctx context.Context, cfg *config.Config) (*mongo.Client, func()) {
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

func connectToRelationalDatabaseAndMigrate(cfg *config.Config) {
	_, err := db.ConnectToDatabaseProvider(cfg.Database, cfg.DatabaseDSN)

	if err != nil {
		log.Fatalf("Error while connection to PostgreSQL: %v", err)
	}

	if err := db.Migrate(cfg.StoreSecretInSql); err != nil {
		log.Fatalf("Auto migration failed: %v", err)
	}
}

func main() {
	cfg, err := config.NewConfig()

	if err != nil {
		log.Fatalf("Error while creating app configuration: %v\n", err)
	}

	go connectToRelationalDatabaseAndMigrate(&cfg)

	ctx, cancel := context.WithCancel(context.Background())
	_, close := connectToMongo(ctx, &cfg)

	close()
	cancel()
}
