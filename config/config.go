package config

import (
	"encoding/base64"
	"errors"
	"os"
	"strconv"
)

type Config struct {
	ApplicationKey []byte
	Database       string
	Mongo          string
	Port           uint16
}

func NewConfig() (Config, error) {
	config := Config{}

	port, err := strconv.ParseUint(os.Getenv("PORT"), 10, 16)
	if err != nil {
		return config, err
	}

	config.Port = uint16(port)

	mongo := os.Getenv("MONGODB")
	if mongo == "" {
		return config, errors.New("MongoDB Url is not set")
	}

	config.Mongo = mongo

	database := os.Getenv("DATABASE_DSN")
	if database == "" {
		return config, errors.New("Datanase DSN is not set")
	}

	config.Database = database

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
