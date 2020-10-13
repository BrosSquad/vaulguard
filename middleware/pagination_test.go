package middleware

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestPaginationMiddleware(t *testing.T) {
	t.Parallel()
	asserts := require.New(t)
	type (
		payload struct {
			Page    int `json:"page"`
			PerPage int `json:"per_page"`
		}
	)

	setup := func() *fiber.App {
		app := fiber.New()
		app.Use(ParsePageAndPerPage)
		app.Get("/", func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusOK).JSON(payload{
				Page:    c.Locals("page").(int),
				PerPage: c.Locals("perPage").(int),
			})
		})
		return app
	}

	t.Run("SuccessfulPageAndPerPageParse", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		query := req.URL.Query()
		query.Set("page", "10")
		query.Set("perPage", "15")
		req.URL.RawQuery = query.Encode()
		app := setup()
		res, err := app.Test(req)
		asserts.Nil(err)
		asserts.Equal(fiber.StatusOK, res.StatusCode)
		p := payload{}
		asserts.Nil(json.NewDecoder(res.Body).Decode(&p))
		asserts.Equal(10, p.Page)
		asserts.Equal(15, p.PerPage)
		asserts.Nil(res.Body.Close())
	})

	t.Run("InvalidPage", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		query := req.URL.Query()
		query.Set("page", "error_page")
		query.Set("perPage", "15")
		req.URL.RawQuery = query.Encode()
		app := setup()
		res, err := app.Test(req)
		asserts.Nil(err)
		asserts.Equal(fiber.StatusUnprocessableEntity, res.StatusCode)
		data, err := ioutil.ReadAll(res.Body)
		asserts.Equal("page query parameter is not a number", string(data))
		asserts.Nil(res.Body.Close())
	})

	t.Run("InvalidPerPage", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		query := req.URL.Query()
		query.Set("page", "10")
		query.Set("perPage", "invalid_perPage")
		req.URL.RawQuery = query.Encode()
		app := setup()
		res, err := app.Test(req)
		asserts.Nil(err)
		asserts.Equal(fiber.StatusUnprocessableEntity, res.StatusCode)
		data, err := ioutil.ReadAll(res.Body)
		asserts.Equal("perPage query parameter is not a number", string(data))
		asserts.Nil(res.Body.Close())
	})

}
