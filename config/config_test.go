package config

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	t.Parallel()
	t.Run("GoodConfiguration", func(t *testing.T) {
		t.Parallel()
		asserts := assert.New(t)
		configStr := `
debug: true
sql: true
http:
  prefork: true
  address: :4000

keys:
  private: ./keys/private
  public: ./keys/public

log:
  level: info

databases:
  mongo:
    uri: mongodb://localhost:27017
  sql:
    provider: postgres
    dsn: 'host=localhost user=postgres pass=postgres dbname=vaulguard timezone=UTC'
`

		buffer := bytes.NewBufferString(configStr)

		config, err := NewConfig(buffer)

		if err != nil {
			t.Fatalf("Error while creating config: %v", err)
		}

		asserts.NotNil(config)
		asserts.Equal(true, config.Debug)
		asserts.Equal(true, config.UseSql)
		asserts.Equal(true, config.Http.Prefork)
		asserts.Equal(":4000", config.Http.Address)
		asserts.Equal("./keys/private", config.Keys.Private)
		asserts.Equal("./keys/public", config.Keys.Public)
		asserts.Equal("info", config.Logging.Level)
		asserts.Equal("mongodb://localhost:27017", config.Databases.Mongo.URI)
		asserts.Equal("host=localhost user=postgres pass=postgres dbname=vaulguard timezone=UTC", config.Databases.SQL.DSN)
		asserts.Equal("postgres", config.Databases.SQL.Provider)

	})

	t.Run("UseSQLDNSEmpty", func(t *testing.T) {
		t.Parallel()
		asserts := assert.New(t)
		configStr := `
debug: true
sql: true
http:
  prefork: true
  address: :4000

keys:
  private: ./keys/private
  public: ./keys/public

log:
  level: info

databases:
  mongo:
    uri: mongodb://localhost:27017
  sql:
    provider:
    dsn:
`
		buffer := bytes.NewBufferString(configStr)

		config, err := NewConfig(buffer)

		asserts.NotNil(err)
		asserts.Nil(config)
		asserts.True(errors.Is(err, ErrDSNEmpty))
	})
	t.Run("UseSQLProviderEmpty", func(t *testing.T) {
		t.Parallel()
		asserts := assert.New(t)
		configStr := `
debug: true
sql: true
http:
  prefork: true
  address: :4000

keys:
  private: ./keys/private
  public: ./keys/public

log:
  level: info

databases:
  mongo:
    uri: mongodb://localhost:27017
  sql:
    provider:
    dsn: 'host=localhost user=postgres pass=postgres dbname=vaulguard timezone=UTC'
`
		buffer := bytes.NewBufferString(configStr)

		config, err := NewConfig(buffer)

		asserts.NotNil(err)
		asserts.Nil(config)
		asserts.True(errors.Is(err, ErrDatabaseProviderEmpty))
	})

	t.Run("MongoURIEmpty", func(t *testing.T) {
		t.Parallel()
		asserts := assert.New(t)
		configStr := `
debug: true
sql: false
http:
  prefork: true
  address: :4000

keys:
  private: ./keys/private
  public: ./keys/public

log:
  level: info

databases:
  mongo:
    uri:
  sql:
    provider: not-supported
    dsn: 'host=localhost user=postgres pass=postgres dbname=vaulguard timezone=UTC'
`
		buffer := bytes.NewBufferString(configStr)

		config, err := NewConfig(buffer)

		asserts.NotNil(err)
		asserts.Nil(config)
		asserts.True(errors.Is(err, ErrMongoURIEmpty))
	})

	t.Run("HttpAddressEmpty", func(t *testing.T) {
		t.Parallel()
		asserts := assert.New(t)
		configStr := `
debug: true
sql: false
http:
  prefork: true
  address:

keys:
  private: ./keys/private
  public: ./keys/public

log:
  level: info

databases:
  mongo:
    uri: mongodb://localhost:27017
  sql:
    provider: postgres
    dsn: 'host=localhost user=postgres pass=postgres dbname=vaulguard timezone=UTC'
`
		buffer := bytes.NewBufferString(configStr)

		config, err := NewConfig(buffer)

		asserts.NotNil(err)
		asserts.Nil(config)
		asserts.True(errors.Is(err, ErrAddressEmpty))
	})

	t.Run("PrivateKeyEmpty", func(t *testing.T) {
		asserts := assert.New(t)
		configStr := `
debug: true
sql: false
http:
  prefork: true
  address: :4000

keys:
  private:
  public: ./keys/public

log:
  level: info

databases:
  mongo:
    uri: mongodb://localhost:27017
  sql:
    provider: postgres
    dsn: 'host=localhost user=postgres pass=postgres dbname=vaulguard timezone=UTC'
`
		buffer := bytes.NewBufferString(configStr)

		config, err := NewConfig(buffer)

		asserts.NotNil(err)
		asserts.Nil(config)
		asserts.True(errors.Is(err, ErrPrivateKeyEmpty))
	})

	t.Run("PrivateKeyEmpty", func(t *testing.T) {
		t.Parallel()
		asserts := assert.New(t)
		configStr := `
debug: true
sql: false
http:
  prefork: true
  address: :4000

keys:
  private: ./keys/private
  public:

log:
  level: info

databases:
  mongo:
    uri: mongodb://localhost:27017
  sql:
    provider: postgres
    dsn: 'host=localhost user=postgres pass=postgres dbname=vaulguard timezone=UTC'
`
		buffer := bytes.NewBufferString(configStr)

		config, err := NewConfig(buffer)

		asserts.NotNil(err)
		asserts.Nil(config)
		asserts.True(errors.Is(err, ErrPublicKeyEmpty))
	})

	t.Run("YamlParseError", func(t *testing.T) {
		t.Parallel()
		asserts := assert.New(t)
		configStr := `
debug: true
sql: 
|false
http:
  prefork: true
  address: :4000

keys:
  private: ./keys/private
  public:

log:
  level: info

databases:
 mongo:
                           uri: mongodb://localhost:27017
  sql:
    provider: postgres
    dsn: 'host=localhost user=postgres pass=postgres dbname=vaulguard timezone=UTC'
`
		buffer := bytes.NewBufferString(configStr)

		config, err := NewConfig(buffer)

		asserts.NotNil(err)
		asserts.Nil(config)
	})

	t.Run("EmptyRead", func(t *testing.T) {
		t.Parallel()
		asserts := assert.New(t)

		config, err := NewConfig(reader{})
		asserts.NotNil(err)
		asserts.Nil(config)
	})
}

type reader struct {}

func (r reader) Read(p []byte) (n int, err error) {
	return 0, io.EOF
}

