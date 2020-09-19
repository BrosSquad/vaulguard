package config

import (
	"errors"
	"github.com/go-yaml/yaml"
	"io"
	"io/ioutil"
	"strings"
)

var (
	ErrDatabaseProviderNotSupported = errors.New("database provider not supported")
	ErrDatabaseProviderEmpty        = errors.New("database provider is required (sqlite, mysql, postgres)")
	ErrDSNEmpty                     = errors.New("DSN is required")
	ErrAddressEmpty                 = errors.New("http address is required")
)

type (
	http struct {
		Prefork bool   `yaml:"prefork,omitempty"`
		Address string `yaml:"address,omitempty"`
	}
	keys struct {
		Private string `yaml:"private,omitempty"`
		Public  string `yaml:"public,omitempty"`
	}
	logging struct {
		Level string `yaml:"level,omitempty"`
	}
	mongo struct {
		URI string `yaml:"uri,omitempty"`
	}
	sql struct {
		Provider string `yaml:"provider,omitempty"`
		DSN      string `yaml:"dsn,omitempty"`
	}

	databases struct {
		Mongo mongo `yaml:"mongo"`
		SQL   sql   `yaml:"sql"`
	}

	Config struct {
		Debug     bool      `yaml:"debug,omitempty"`
		UseSql    bool      `yaml:"sql,omitempty"`
		Http      http      `yaml:"http"`
		Keys      keys      `yaml:"keys"`
		Logging   logging   `yaml:"log"`
		Databases databases `yaml:"databases"`
	}
)

func checkDatabaseProvider(provider string) error {
	providers := [3]string{"postgres", "mysql", "sqlite"}

	provider = strings.ToLower(provider)

	for _, p := range providers {
		if p == provider {
			return nil
		}
	}

	return ErrDatabaseProviderNotSupported
}

func (c Config) Validate() error {

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

	return &config, nil
}
