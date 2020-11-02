package config

import (
	"errors"
	"github.com/go-yaml/yaml"
	"io"
	"time"
)

var (
	ErrDatabaseProviderEmpty = errors.New("database provider is required")
	ErrDSNEmpty              = errors.New("DSN is required")
	ErrMongoURIEmpty         = errors.New("mongo URI is required")
	ErrAddressEmpty          = errors.New("http address is required")
	ErrPrivateKeyEmpty       = errors.New("private key is required")
	ErrPublicKeyEmpty        = errors.New("public key is required")
	ErrLocaleNotFound        = errors.New("locale is required for validation")
	ErrMemoryUsageSleepEmpty = errors.New("memory usage sleep is required")
)

type (
	Http struct {
		Prefork bool   `yaml:"prefork,omitempty"`
		Address string `yaml:"address,omitempty"`
	}
	Keys struct {
		Private string `yaml:"private,omitempty"`
		Public  string `yaml:"public,omitempty"`
		Secret  string `yaml:"secret,omitempty"`
	}
	Logging struct {
		Level string `yaml:"level,omitempty"`
	}
	Mongo struct {
		URI string `yaml:"uri,omitempty"`
	}
	Sql struct {
		Provider string `yaml:"provider,omitempty"`
		DSN      string `yaml:"dsn,omitempty"`
	}

	Databases struct {
		Mongo Mongo `yaml:"mongo,omitempty"`
		SQL   Sql   `yaml:"sql,omitempty"`
	}

	MemoryUsage struct {
		Report bool          `yaml:"report,omitempty"`
		Sleep  time.Duration `yaml:"sleep,omitempty"`
	}

	Config struct {
		ApplicationKey []byte      `yaml:"-"`
		Locale         string      `yaml:"locale,omitempty"`
		UseConsole     bool        `yaml:"console,omitempty"`
		Debug          bool        `yaml:"debug,omitempty"`
		UseSql         bool        `yaml:"sql,omitempty"`
		Http           Http        `yaml:"http,omitempty"`
		Keys           Keys        `yaml:"keys,omitempty"`
		Logging        Logging     `yaml:"log,omitempty"`
		Databases      Databases   `yaml:"databases,omitempty"`
		MemoryUsage    MemoryUsage `yaml:"memory,omitempty"`
	}
)

func (c Config) Validate() error {
	if c.UseSql {
		if c.Databases.SQL.DSN == "" {
			return ErrDSNEmpty
		}
		if c.Databases.SQL.Provider == "" {
			return ErrDatabaseProviderEmpty
		}
	} else {
		if c.Databases.Mongo.URI == "" {
			return ErrMongoURIEmpty
		}
	}

	if c.Http.Address == "" {
		return ErrAddressEmpty
	}

	if c.Keys.Private == "" {
		return ErrPrivateKeyEmpty
	}

	if c.Keys.Public == "" {
		return ErrPublicKeyEmpty
	}

	if c.Locale == "" {
		return ErrLocaleNotFound
	}

	if c.MemoryUsage.Report && c.MemoryUsage.Sleep == 0 {
		return ErrMemoryUsageSleepEmpty
	}
	return nil
}

func NewConfig(r io.Reader) (*Config, error) {
	config := Config{}
	if err := yaml.NewDecoder(r).Decode(&config); err != nil {
		return nil, err
	}
	if err := config.Validate(); err != nil {
		return nil, err
	}
	return &config, nil
}
