package middleware

import (
	"github.com/BrosSquad/vaulguard/models"
	"github.com/BrosSquad/vaulguard/services/token"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type mockTokenService struct {
	mock.Mock
}

func (m *mockTokenService) Generate(i interface{}) string {
	args := m.Called(i)
	return args.String(0)
}

func (m *mockTokenService) Verify(s string) (models.ApplicationDto, bool) {
	args := m.Called(s)

	return args.Get(0).(models.ApplicationDto), args.Bool(1)
}

func TestTokenAuth(t *testing.T) {
	t.Parallel()
	asserts := require.New(t)

	t.Run("NotEnoughParams", func(t *testing.T) {
		asserts.Panics(func() {
			TokenAuth(TokenAuthConfig{
				Headers:        []string{"authorization"},
				HeaderPrefixes: []string{"test"},
				TokenServices:  []token.Service{},
			})
		})
	})

	t.Run("TokenVerificationMiddleware", func(t *testing.T) {
		app := fiber.New()
		mockService := &mockTokenService{}

		mockService.On("Verify", "Test.1.TestToken").Return(models.ApplicationDto{
			ID:        1,
			Name:      "TestApplication",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}, true)

		app.Use(TokenAuth(TokenAuthConfig{
			Headers:        []string{"authorization"},
			HeaderPrefixes: []string{"test "},
			TokenServices:  []token.Service{mockService},
		}))

		app.Get("/", func(ctx *fiber.Ctx) error {
			application := ctx.Locals("application").(models.ApplicationDto)
			_, err := ctx.WriteString("Hello " + application.Name)
			return err
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Add("Authorization", "test Test.1.TestToken")

		resp, err := app.Test(req)
		asserts.Nil(err)
		data, err := ioutil.ReadAll(resp.Body)
		asserts.Nil(err)
		asserts.Equal("Hello TestApplication", string(data))
		asserts.Nil(resp.Body.Close())
	})
	t.Run("TokenVerificationMiddleware_NoToken", func(t *testing.T) {
		app := fiber.New()
		mockService := &mockTokenService{}

		mockService.On("Verify", "Test.1.TestToken").Return(models.ApplicationDto{
			ID:        1,
			Name:      "TestApplication",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}, true)

		app.Use(TokenAuth(TokenAuthConfig{
			Headers:        []string{"authorization"},
			HeaderPrefixes: []string{"test "},
			TokenServices:  []token.Service{mockService},
		}))

		app.Get("/", func(ctx *fiber.Ctx) error {
			application := ctx.Locals("application").(models.ApplicationDto)
			_, err := ctx.WriteString("Hello " + application.Name)
			return err
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)

		resp, err := app.Test(req)
		asserts.Nil(err)
		asserts.Equal(resp.StatusCode, fiber.StatusUnauthorized)
		asserts.Nil(resp.Body.Close())
	})

	t.Run("TokenVerificationMiddleware_InvalidToken", func(t *testing.T) {
		app := fiber.New()
		mockService := &mockTokenService{}

		mockService.On("Verify", "Test.1.TestToken").Return(models.ApplicationDto{}, false)

		app.Use(TokenAuth(TokenAuthConfig{
			Headers:        []string{"authorization"},
			HeaderPrefixes: []string{"test "},
			TokenServices:  []token.Service{mockService},
		}))

		app.Get("/", func(ctx *fiber.Ctx) error {
			application := ctx.Locals("application").(models.ApplicationDto)
			_, err := ctx.WriteString("Hello " + application.Name)
			return err
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Add("Authorization", "test Test.1.TestToken")

		resp, err := app.Test(req)
		asserts.Nil(err)
		asserts.Equal(resp.StatusCode, fiber.StatusUnauthorized)
		asserts.Nil(resp.Body.Close())
	})
}
