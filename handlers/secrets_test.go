package handlers

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"errors"
	"github.com/BrosSquad/vaulguard/services"
	"github.com/BrosSquad/vaulguard/services/application"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/BrosSquad/vaulguard/models"
	"github.com/BrosSquad/vaulguard/services/secret"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockSecretService struct {
	Id      uint
	Mutex   *sync.RWMutex
	IdMutex *sync.Mutex
	Data    []models.Secret
	mock.Mock
}

func (m *mockSecretService) Paginate(applicationID interface{}, page, perPage int) (map[string]string, error) {
	panic("implement me")
}

func (m *mockSecretService) Get(applicationID interface{}, key []string) (map[string]string, error) {
	panic("implement me")
}

func (m *mockSecretService) GetOne(applicationID interface{}, key string) (secret.Secret, error) {
	panic("implement me")
}

func (m *mockSecretService) Create(applicationID interface{}, key, value string) (models.Secret, error) {
	args := m.Called(applicationID, key, value)

	if err := args.Error(0); err != nil {
		return models.Secret{}, err
	}
	m.Mutex.Lock()
	m.IdMutex.Lock()
	defer m.IdMutex.Unlock()
	defer m.Mutex.Unlock()
	m.Id++
	s := models.Secret{ID: m.Id, Key: key, Value: []byte(value), ApplicationId: applicationID.(uint)}
	m.Data = append(m.Data, s)

	return s, nil
}

func (m *mockSecretService) Update(applicationID interface{}, key, newKey, value string) (models.Secret, error) {
	panic("implement me")
}

func (m *mockSecretService) Delete(applicationID interface{}, key string) error {
	panic("implement me")
}

func (m *mockSecretService) InvalidateCache(applicationID interface{}) error {
	panic("implement me")
}

func TestCreateSecret(t *testing.T) {
	t.Parallel()
	asserts := require.New(t)
	v := validator.New()
	english := en.New()
	uni := ut.New(english, english)
	englishTranslations, found := uni.GetTranslator("en")

	asserts.True(found)

	createMockService := func() *mockSecretService {
		return &mockSecretService{
			Id:      0,
			Mutex:   &sync.RWMutex{},
			IdMutex: &sync.Mutex{},
			Data:    make([]models.Secret, 0),
		}
	}

	setup := func(service secret.Service, setupMiddleware bool) *fiber.App {
		app := fiber.New(fiber.Config{
			ErrorHandler: Error(englishTranslations),
		})
		if setupMiddleware {
			app.Use(func(c *fiber.Ctx) error {
				c.Locals("application", models.ApplicationDto{
					ID:        uint(1),
					Name:      "Test Application",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				})
				return c.Next()
			})
			RegisterSecretHandlers(v, service, app.Group("/secrets"))
		}
		return app
	}

	t.Run("InsertSuccess", func(t *testing.T) {
		service := createMockService()
		service.On("Create", uint(1), "Test", "Test").Return(nil)
		app := setup(service, true)

		data, err := json.Marshal(struct {
			Key   string
			Value string
		}{Key: "Test", Value: "Test"})
		asserts.Nil(err)
		buff := bytes.NewBuffer(data)

		req := httptest.NewRequest(http.MethodPost, "/secrets", buff)
		req.Header.Add(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)

		res, err := app.Test(req, 400)
		asserts.Nil(err)
		asserts.EqualValues(http.StatusCreated, res.StatusCode)
		asserts.Len(service.Data, 1)
		asserts.EqualValues(1, service.Id)
	})

	t.Run("InsertFailed", func(t *testing.T) {
		service := createMockService()
		service.On("Create", uint(1), "Test", "Test").Return(errors.New("insert error"))
		app := setup(service, true)

		data, err := json.Marshal(struct {
			Key   string
			Value string
		}{Key: "Test", Value: "Test"})
		asserts.Nil(err)
		buff := bytes.NewBuffer(data)

		req := httptest.NewRequest(http.MethodPost, "/secrets", buff)
		req.Header.Add(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)

		res, err := app.Test(req, 400)
		asserts.Nil(err)
		asserts.EqualValues(http.StatusInternalServerError, res.StatusCode)
		asserts.Len(service.Data, 0)
		asserts.EqualValues(0, service.Id)
	})

	t.Run("ValidationError", func(t *testing.T) {
		service := createMockService()
		app := setup(service, true)

		data, err := json.Marshal(struct {
			Key   string
			Value string
		}{Key: "", Value: ""})
		asserts.Nil(err)
		buff := bytes.NewBuffer(data)

		req := httptest.NewRequest(http.MethodPost, "/secrets", buff)
		req.Header.Add(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)

		res, err := app.Test(req, 400)
		asserts.Nil(err)
		asserts.EqualValues(http.StatusUnprocessableEntity, res.StatusCode)
		asserts.Len(service.Data, 0)
		asserts.EqualValues(0, service.Id)
	})

	t.Run("GormIntegrationTest", func(t *testing.T) {
		asserts := require.New(t)
		key := make([]byte, services.SecretKeyLength)
		n, err := rand.Read(key)
		asserts.Nil(err)
		asserts.Equal(services.SecretKeyLength, n)
		encryption, err := services.NewSecretKeyEncryption(key)
		asserts.Nil(err)
		path, err := filepath.Abs("./create_secrets.db")
		asserts.Nil(err)
		defer os.Remove(path)
		db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
		asserts.Nil(err)
		asserts.Nil(db.AutoMigrate(&models.Application{},&models.Token{},&models.Secret{}))
		service := secret.NewGormSecretStorage(secret.GormSecretConfig{
			Encryption: encryption,
			CacheSize:  10,
			DB:         db,
		})
		applicationDto, err := application.NewSqlService(db).Create("TestApplication")
		asserts.Nil(err)
		app := setup(service, false)
		app.Use(func(c *fiber.Ctx) error {
			c.Locals("application", applicationDto)
			return c.Next()
		})
		RegisterSecretHandlers(v, service, app.Group("/secrets"))
		data, err := json.Marshal(struct {
			Key   string
			Value string
		}{Key: "Test", Value: "Test"})
		asserts.Nil(err)
		buff := bytes.NewBuffer(data)

		req := httptest.NewRequest(http.MethodPost, "/secrets", buff)
		req.Header.Add(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)

		res, err := app.Test(req, 400)
		asserts.Nil(err)
		asserts.Equal(fiber.StatusCreated, res.StatusCode)
		payload := struct {
			ID uint `json:"id"`
			Key string `json:"key"`
			Value string `json:"value"`
		}{}
		err = json.NewDecoder(res.Body).Decode(&payload)
		asserts.Nil(err)
		asserts.NotEqual(0, payload.ID)
		asserts.Equal("Test", payload.Key)
		asserts.Equal("Test", payload.Value)
	})

}
