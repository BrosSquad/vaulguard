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
	Debug          bool
	ApplicationKey []byte
	StoreInSql     bool
	Database       string
	DatabaseDSN    string
	Mongo          string
	Port           string
}

var (
	errSkipEnv = errors.New("Skip Enviromental variables.")
	ErrDatabaseProviderEmpty = errors.New("SQL Database provider is required.")
	ErrDSNEmpty = errors.New("Database DSN(Data Source Name) is required.")
	ErrDatabaseProviderNotSupported = errors.New("SQL Database provider is not supported.")
	ErrAppKeyEmpty = errors.New("App Key is required.")
	ErrAppKeyLength = errors.New("App key has to be 32 bytes in length.")
	ErrAddressEmpty =  errors.New("Address is required.")
	ErrMongoURIEmpty = errors.New("MongoURI is required.")
)

func checkDatabaseProvider(provider string) error {
	providers := []string{"postgres", "mysql", "sqlite"}

	provider = strings.ToLower(provider)

	for _, p := range providers {
		if p == provider {
			return nil
		}
	}

	return ErrDatabaseProviderNotSupported
}

func handleFlags(cfg *Config) error {
	debug := flag.Bool("debug", false, "Debug mode - Command line arguments are only accepted in DEBUG mode")
	skipEnv := flag.Bool("skip-env", false, "Skips Environmental variables and only uses command line flags")
	port := flag.String("port", ":4000", "HTTP Server port")
	storeSecretsInSQL := flag.Bool("store-in-sql", false, "Store data in SQL database")
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
	cfg.StoreInSql = *storeSecretsInSQL
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

func (c *Config) handleSqlDatabase() error {
	c.Database = os.Getenv("VAULGUARD_DATABASE_PROVIDER")
	if c.Database == "" {
		return ErrDatabaseProviderEmpty
	}

	if err := checkDatabaseProvider(c.Database); err != nil {
		return err
	}

	c.DatabaseDSN = os.Getenv("VAULGUARD_DATABASE_DSN")

	if c.DatabaseDSN == "" {
		return ErrDSNEmpty
	}

	return nil
}

func (c *Config) handleAppKey() error {
	appKey := os.Getenv("VAULGUARD_APP_KEY")

	if appKey == "" {
		return ErrAppKeyEmpty
	}

	appKeyBytes, err := base64.StdEncoding.DecodeString(appKey)

	if err != nil {
		return err
	}

	if len(appKeyBytes) != 32 {
		return ErrAppKeyLength
	}

	c.ApplicationKey = appKeyBytes

	return nil
}

func (c *Config) handlePort() error {
	c.Port = os.Getenv("VAULGUARD_ADDRESS")

	if c.Port == "" {
		return ErrAddressEmpty
	}

	return nil
}

func (c *Config) handleDatabase() error {
	storeInSQL, err := strconv.ParseBool(os.Getenv("VAULGUARD_STORE_IN_SQL"))

	if err == nil {
		c.StoreInSql = storeInSQL
	}

	if !storeInSQL {
		c.Mongo = os.Getenv("VAULGUARD_MONGODB")
		if c.Mongo == "" {
			return ErrMongoURIEmpty
		}
	} else {
		if err := c.handleSqlDatabase(); err != nil {
			return err
		}
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

	if err := config.handlePort(); err != nil {
		return Config{}, err
	}

	if err := config.handleDatabase(); err != nil {
		return Config{}, err
	}

	if err := config.handleAppKey(); err != nil {
		return Config{}, err
	}

	return config, nil
}
