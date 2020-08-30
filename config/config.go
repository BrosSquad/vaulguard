package config

import (
	"encoding/base64"
	"errors"
	"flag"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Debug            bool
	ApplicationKey   []byte
	StoreSecretInSql bool
	Database         string
	DatabaseDSN      string
	Mongo            string
	Port             int
}

var errSkipEnv = errors.New("Skip Enviromental variables")

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

func handleFlags(cfg *Config) error {
	debug := flag.Bool("debug", false, "Debug mode - Command line arguments are only accepted in DEBUG mode")
	skipEnv := flag.Bool("skip-env", false, "Skips Environmental variables and only uses command line flags")
	port := flag.Int("port", 4000, "HTTP Server port")
	storeSecretsInSQL := flag.Bool("store-secrets-in-sql", false, "Store secrets in SQL database")
	mongo := flag.String("mongo", "", "MongoDB connection string")
	dbProvider := flag.String("provider", "sqlite", "Relational database provider (sqlite, mysql, postgres)")
	connectionString := flag.String("db-connection", "", "Relational database connection string")
	appKey := flag.String("app-key", "", "Application encryption key")
	flag.Parse()

	if !*debug {
		return nil
	}

	cfg.Debug = *debug
	cfg.Port = *port
	cfg.StoreSecretInSql = *storeSecretsInSQL
	cfg.Mongo = *mongo
	cfg.Database = *dbProvider
	cfg.DatabaseDSN = *connectionString

	if *appKey != "" {
		key, err := base64.StdEncoding.DecodeString(*appKey)
		if err != nil {
			return err
		}
		cfg.ApplicationKey = key
	}

	if *skipEnv {
		return errSkipEnv
	}

	return nil
}

func NewConfig(skipFlags ...bool) (Config, error) {
	config := Config{}

	if len(skipFlags) > 0 && skipFlags[0] == true {
		err := handleFlags(&config)

		if errors.Is(err, errSkipEnv) {
			return config, nil
		}

		if err != nil {
			return config, err
		}
	}

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
