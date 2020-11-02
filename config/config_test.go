package config

import (
	"bytes"
	"errors"
	"io"
	"testing"
	"time"

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
locale: en
http:
  prefork: true
  address: :4000
memory:
  report: true
  sleep: 40s
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

		config, err := New(buffer)

		if err != nil {
			t.Fatalf("Error while creating config: %v", err)
		}

		asserts.NotNil(config)
		asserts.True(config.Debug)
		asserts.True(config.UseSql)
		asserts.True(config.Http.Prefork)
		asserts.EqualValues(time.Duration(40)*time.Second, config.MemoryUsage.Sleep)
		asserts.True(config.MemoryUsage.Report)
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
locale: en
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

		config, err := New(buffer)

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
locale: en
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

		config, err := New(buffer)

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
locale: en
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

		config, err := New(buffer)

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
locale: en
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

		config, err := New(buffer)

		asserts.NotNil(err)
		asserts.Nil(config)
		asserts.True(errors.Is(err, ErrAddressEmpty))
	})

	t.Run("PrivateKeyEmpty", func(t *testing.T) {
		asserts := assert.New(t)
		configStr := `
debug: true
locale: en
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

		config, err := New(buffer)

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
locale: en
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

		config, err := New(buffer)

		asserts.NotNil(err)
		asserts.Nil(config)
		asserts.True(errors.Is(err, ErrPublicKeyEmpty))
	})

	t.Run("YamlParseError", func(t *testing.T) {
		t.Parallel()
		asserts := assert.New(t)
		configStr := `
debug: true
locale: en
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

		config, err := New(buffer)

		asserts.NotNil(err)
		asserts.Nil(config)
	})

	t.Run("EmptyRead", func(t *testing.T) {
		t.Parallel()
		asserts := assert.New(t)

		config, err := New(reader{})
		asserts.NotNil(err)
		asserts.Nil(config)
	})

	t.Run("LocaleEmpty", func(t *testing.T) {
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
  secret: ./keys/secret
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

		config, err := New(buffer)

		asserts.NotNil(err)
		asserts.Nil(config)
		asserts.True(errors.Is(err, ErrLocaleNotFound))
	})
	t.Run("MemoryUsageSleepIsZero", func(t *testing.T) {
		t.Parallel()
		asserts := assert.New(t)
		configStr := `
debug: true
sql: false
locale: en
http:
  prefork: true
  address: :4000
memory:
  report: true
  sleep: 0
keys:
  private: ./keys/private
  public: ./keys/public
  secret: ./keys/secret
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

		config, err := New(buffer)

		asserts.NotNil(err)
		asserts.Nil(config)
		asserts.True(errors.Is(err, ErrMemoryUsageSleepEmpty))
	})
}

type reader struct{}

func (r reader) Read(p []byte) (n int, err error) {
	return 0, io.EOF
}
