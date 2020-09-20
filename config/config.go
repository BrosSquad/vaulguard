package config

import (
	"errors"
	"github.com/go-yaml/yaml"
	"io"
	"io/ioutil"
)

var (
	ErrDatabaseProviderEmpty = errors.New("database provider is required")
	ErrDSNEmpty              = errors.New("DSN is required")
	ErrMongoURIEmpty         = errors.New("mongo URI is required")
	ErrAddressEmpty          = errors.New("http address is required")
	ErrPrivateKeyEmpty       = errors.New("private key is required")
	ErrPublicKeyEmpty        = errors.New("public key is required")
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

	databases struct {
		Mongo Mongo `yaml:"mongo"`
		SQL   Sql   `yaml:"sql"`
	}

	Config struct {
		ApplicationKey []byte    `yaml:"-"`
		UseConsole     bool      `yaml:"console,omitempty"`
		Debug          bool      `yaml:"debug,omitempty"`
		UseSql         bool      `yaml:"sql,omitempty"`
		Http           Http      `yaml:"http"`
		Keys           Keys      `yaml:"keys"`
		Logging        Logging   `yaml:"log"`
		Databases      databases `yaml:"databases"`
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

	return nil
}

func NewConfig(r io.Reader) (*Config, error) {
	bytes, err := ioutil.ReadAll(r)

	if err != nil {
		return nil, err
	}

	config := Config{}
	if err := yaml.Unmarshal(bytes, &config); err != nil {
		return nil, err
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &config, nil
}
