package config

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/go-yaml/yaml"
)

const EnvironmentalVariablesPrefix = "VAULGUARD_"

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
	Session struct {
		CookieName string        `yaml:"cookie,omitempty"`
		Provider   string        `yaml:"provider,omitempty"`
		Domain     string        `yaml:"domain,omitempty"`
		SameSite   string        `yaml:"samesite,omitempty"`
		Expiration time.Duration `yaml:"expiration,omitempty"`
		GC         time.Duration `yaml:"gc,omitempty"`
		RedisDB    int64         `yaml:"redi_db,omitempty"`
		Secure     bool          `yaml:"secure,omitempty"`
	}

	Http struct {
		Address string  `yaml:"address,omitempty"`
		Session Session `yaml:"session:omitempty"`
		Prefork bool    `yaml:"prefork,omitempty"`
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

	Redis struct {
		Addr     string `yaml:"addr,omitempty"`
		Password string `yaml:"password,omitempty"`
	}

	Databases struct {
		Mongo Mongo `yaml:"mongo,omitempty"`
		SQL   Sql   `yaml:"sql,omitempty"`
		Redis Redis `yaml:"redis,omitempty"`
	}

	MemoryUsage struct {
		Report bool          `yaml:"report,omitempty"`
		Sleep  time.Duration `yaml:"sleep,omitempty"`
	}

	Config struct {
		ApplicationKey []byte      `yaml:"-"`
		Locale         string      `yaml:"locale,omitempty"`
		Http           Http        `yaml:"http,omitempty"`
		Keys           Keys        `yaml:"keys,omitempty"`
		Logging        Logging     `yaml:"log,omitempty"`
		Databases      Databases   `yaml:"databases,omitempty"`
		MemoryUsage    MemoryUsage `yaml:"memory,omitempty"`
		UseConsole     bool        `yaml:"console,omitempty"`
		Debug          bool        `yaml:"debug,omitempty"`
		UseSql         bool        `yaml:"sql,omitempty"`
		UseDashboard   bool        `yaml:"dashboard,omitempty"`
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

	// TODO: Add HTTP Session Validation
	if c.UseDashboard {

	}

	return nil
}

func (c *Config) loadEnvironmentalVariables() (err error) {
	prefork := os.Getenv(EnvironmentalVariablesPrefix + "HTTP_PREFORK")
	if prefork != "" {
		c.Http.Prefork, err = strconv.ParseBool(prefork)
		if err != nil {
			return err
		}
	}

	address := os.Getenv(EnvironmentalVariablesPrefix + "HTTP_ADDRESS")
	if address != "" {
		c.Http.Address = address
	}

	publicKey := os.Getenv(EnvironmentalVariablesPrefix + "PUBLIC_KEY")
	if publicKey != "" {
		c.Keys.Public = publicKey
	}

	privateKey := os.Getenv(EnvironmentalVariablesPrefix + "PRIVATE_KEY")
	if privateKey != "" {
		c.Keys.Public = privateKey
	}

	secretKey := os.Getenv(EnvironmentalVariablesPrefix + "SECRET_KEY")
	if secretKey != "" {
		c.Keys.Public = secretKey
	}

	loggingLevel := os.Getenv(EnvironmentalVariablesPrefix + "LOGGING_LEVEL")
	if loggingLevel != "" {
		c.Logging.Level = loggingLevel
	}

	mongo := os.Getenv(EnvironmentalVariablesPrefix + "MONGO_URI")
	if mongo != "" {
		c.Databases.Mongo.URI = mongo
	}

	sqlProvider := os.Getenv(EnvironmentalVariablesPrefix + "SQL_PROVIDER")
	if sqlProvider != "" {
		c.Databases.SQL.Provider = sqlProvider
	}

	sqlDSN := os.Getenv(EnvironmentalVariablesPrefix + "SQL_DSN")
	if sqlDSN != "" {
		c.Databases.SQL.Provider = sqlDSN
	}

	locale := os.Getenv(EnvironmentalVariablesPrefix + "LOCALE")
	if locale != "" {
		c.Locale = locale
	}

	useConsole := os.Getenv(EnvironmentalVariablesPrefix + "USE_CONSOLE")
	if useConsole != "" {
		c.UseConsole, err = strconv.ParseBool(useConsole)
		if err != nil {
			return err
		}
	}

	debug := os.Getenv(EnvironmentalVariablesPrefix + "DEBUG")
	if debug != "" {
		c.Debug, err = strconv.ParseBool(debug)
		if err != nil {
			return err
		}
	}

	useSQL := os.Getenv(EnvironmentalVariablesPrefix + "USE_SQL")
	if useSQL != "" {
		c.UseSql, err = strconv.ParseBool(useSQL)
		if err != nil {
			return err
		}
	}

	memoryUsageReport := os.Getenv(EnvironmentalVariablesPrefix + "MEMORY_USAGE_REPORT")
	if memoryUsageReport != "" {
		c.MemoryUsage.Report, err = strconv.ParseBool(memoryUsageReport)
		if err != nil {
			return err
		}
	}

	memoryUsageSleep := os.Getenv(EnvironmentalVariablesPrefix + "MEMORY_USAGE_SLEEP")
	if memoryUsageSleep != "" {
		c.MemoryUsage.Sleep, err = time.ParseDuration(memoryUsageSleep)
		if err != nil {
			return err
		}
	}

	sessionCookieName := os.Getenv(EnvironmentalVariablesPrefix + "SESSION_COOKIE_NAME")
	if sessionCookieName != "" {
		c.Http.Session.CookieName = sessionCookieName
	}

	sessionStorageProvider := os.Getenv(EnvironmentalVariablesPrefix + "SESSION_PROVIDER")
	if sessionStorageProvider != "" {
		c.Http.Session.CookieName = sessionStorageProvider
	}

	sessionSecure := os.Getenv(EnvironmentalVariablesPrefix + "SESSION_SECURE")
	if sessionSecure != "" {
		c.Http.Session.Secure, err = strconv.ParseBool(sessionSecure)
		if err != nil {
			return err
		}
	}

	sessionDomain := os.Getenv(EnvironmentalVariablesPrefix + "SESSION_DOMAIN")
	if sessionDomain != "" {
		c.Http.Session.Domain = sessionDomain
	}

	sessionSameSite := os.Getenv(EnvironmentalVariablesPrefix + "SESSION_SAMESITE")
	if sessionSameSite != "" {
		c.Http.Session.SameSite = sessionSameSite
	}

	sessionExpiration := os.Getenv(EnvironmentalVariablesPrefix + "MEMORY_EXPIRATION")
	if sessionExpiration != "" {
		c.Http.Session.Expiration, err = time.ParseDuration(sessionExpiration)
		if err != nil {
			return err
		}
	}

	sessionGC := os.Getenv(EnvironmentalVariablesPrefix + "MEMORY_GC")
	if sessionGC != "" {
		c.Http.Session.GC, err = time.ParseDuration(sessionGC)
		if err != nil {
			return err
		}
	}

	sessionRedisDB := os.Getenv(EnvironmentalVariablesPrefix + "REDIS_DB")
	if sessionRedisDB != "" {
		c.Http.Session.RedisDB, err = strconv.ParseInt(sessionRedisDB, 10, 32)
		if err != nil {
			return err
		}
	}

	redisAddr := os.Getenv(EnvironmentalVariablesPrefix + "REDIS_ADDR")
	if redisAddr != "" {
		c.Databases.Redis.Addr = redisAddr
	}

	redisPassword := os.Getenv(EnvironmentalVariablesPrefix + "REDIS_PASSWORD")
	if redisPassword != "" {
		c.Databases.Redis.Addr = redisPassword
	}

	return nil
}

func New(r io.Reader) (*Config, error) {
	bytes, err := ioutil.ReadAll(r)

	if err != nil {
		return nil, err
	}

	config := Config{}
	if err := yaml.Unmarshal(bytes, &config); err != nil {
		return nil, err
	}

	if err := config.loadEnvironmentalVariables(); err != nil {
		return nil, err
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &config, nil
}
