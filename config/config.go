package config

import (
	"encoding/base64"
	"errors"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	ApplicationKey   []byte
	StoreSecretInSql bool
	Database         string
	DatabaseDSN      string
	Mongo            string
	Port             int
}

func checkDatabaseProvider(provider string) error {
	providers := []string{"postgres", "mysql", "sqlite"}

	provider = strings.ToLower(provider)

	for _, p := range providers {
		if p == provider {
			return nil
		}
	}

	return errors.New("Database provider is not supported")
}

func NewConfig() (Config, error) {
	config := Config{}

	port, err := strconv.ParseInt(os.Getenv("PORT"), 10, 16)

	if err != nil {
		return config, err
	}

	config.Port = int(port)

	storeInSQL, err := strconv.ParseBool(os.Getenv("STORE_SECRETS_IN_SQL"))

	if err == nil {
		config.StoreSecretInSql = storeInSQL
	}

	if storeInSQL {
		mongo := os.Getenv("MONGODB")
		if mongo == "" {
			return config, errors.New("MongoDB Url is not set")
		}
		config.Mongo = mongo
	}

	database := os.Getenv("DATABASE_PROVIDER")
	if database == "" {
		return config, errors.New("Database provider is not set")
	}

	if err := checkDatabaseProvider(database); err != nil {
		return config, err
	}
	config.Database = database

	databaseDsn := os.Getenv("DATABASE_DSN")
	if databaseDsn == "" {
		return config, errors.New("Database DSN is not set")
	}

	config.DatabaseDSN = databaseDsn

	appKey := os.Getenv("APP_KEY")

	if appKey == "" {
		return config, errors.New("Application Key is not set")
	}

	appKeyBytes, err := base64.StdEncoding.DecodeString(appKey)

	if err != nil {
		return config, err
	}

	if len(appKeyBytes) != 32 {
		return config, errors.New("Application key has to be 32 bytes long")
	}

	config.ApplicationKey = appKeyBytes

	return config, nil
}
