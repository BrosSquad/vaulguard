package main

import (
	"context"
	"log"
	"time"

	"github.com/BrosSquad/vaulguard/config"
	"github.com/BrosSquad/vaulguard/db"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	cfg, err := config.NewConfig()

	if err != nil {
		log.Fatalf("Error while creating app configuration: %v\n", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := db.ConnectToMongo(ctx, cfg.Mongo)

	if err != nil {
		log.Fatalf("Error while connecting to mongo db instance: %v\n", err)
	}

	log.Println(client.Ping(ctx, readpref.Primary()))

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatalf("Error while disconecting from mongodb instance: %v\n", err)
		}
	}()

	_, err = db.ConnectToDatabaseProvider(cfg.Database, cfg.DatabaseDSN)

	if err != nil {
		log.Fatalf("Error while connection to PostgreSQL: %v", err)
	}

}
