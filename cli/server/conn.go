package main

import (
	"context"
	"io"

	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"

	"github.com/BrosSquad/vaulguard/config"
	"github.com/BrosSquad/vaulguard/db"
	"github.com/BrosSquad/vaulguard/log"
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

func connectToMongoAndMigrate(ctx context.Context, mongoURI string) (*mongo.Client, io.Closer, error) {
	client, err := db.ConnectToMongo(ctx, mongoURI)

	if err != nil {
		return nil, nil, err
	}

	if err := db.MongoCreateCollections(ctx, client); err != nil {
		return nil, nil, err
	}

	return client, &mongoClose{ctx, client}, nil
}

func connectToRelationalDatabaseAndMigrate(logger *log.Logger, cfg *config.Config) (*gorm.DB, io.Closer, error) {
	provider, err := db.GetDatabaseProvider(cfg.Databases.SQL.Provider)
	if err != nil {
		return nil, nil, err
	}

	conn, err := db.ConnectToDatabaseProvider(db.GormConfig{
		LogLevel:    log.GetDbLogLevel(cfg.Logging.Level),
		Logger:      logger,
		SQLProvider: provider,
		DSN:         cfg.Databases.SQL.DSN,
	})

	if err != nil {
		return nil, nil, err
	}

	if err := db.SqlMigrate(); err != nil {
		return nil, nil, err
	}

	return conn, &gormClose{conn}, nil
}
