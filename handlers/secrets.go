package handlers

import (
	"github.com/BrosSquad/vaulguard/middleware"
	"github.com/BrosSquad/vaulguard/models"
	"github.com/BrosSquad/vaulguard/services/secret"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func RegisterSecretHandlers(validate *validator.Validate, service secret.Service, r fiber.Router) {
	r.Get("/", getSecrets(service)).Use(middleware.ParsePageAndPerPage)
	r.Get("/many", getManySecrets(service))
	r.Post("/", createSecret(service, validate))
	r.Delete("/invalidate", invalidateCache(service))
}

func getSecrets(service secret.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		app := ctx.Locals("application").(models.ApplicationDto)
		page := ctx.Locals("page").(int)
		perPage := ctx.Locals("perPage").(int)

		secrets, err := service.Paginate(app.ID, page, perPage)

		if err != nil {
			return err
		}

		return ctx.JSON(secrets)
	}
}

func getManySecrets(service secret.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		keysStruct := struct {
			keys []string `query:"keys"`
		}{}
		app := c.Locals("application").(models.ApplicationDto)
		if err := c.QueryParser(&keysStruct); err != nil {
			return fiber.ErrBadRequest
		}

		secrets, err := service.Get(app.ID, keysStruct.keys)
		if err != nil {
			return err
		}

		return c.JSON(secrets)
	}
}

func createSecret(service secret.Service, validate *validator.Validate) fiber.Handler {
	return func(c *fiber.Ctx) error {
		payload := struct {
			Key   string `json:"key" validate:"required"`
			Value string `json:"value" validate:"required"`
		}{}

		app := c.Locals("application").(models.ApplicationDto)
		if err := c.BodyParser(&payload); err != nil {
			return fiber.ErrBadRequest
		}

		if err := validate.Struct(payload); err != nil {
			return err
		}

		s, err := service.Create(app.ID, payload.Key, payload.Value)

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusCreated).JSON(struct {
			ID    interface{} `json:"id"`
			Key   string      `json:"key"`
			Value string      `json:"value"`
		}{
			ID:    s.ID,
			Key:   s.Key,
			Value: payload.Value,
		})
	}
}

func invalidateCache(service secret.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		app := ctx.Locals("application").(models.ApplicationDto)

		if err := service.InvalidateCache(app.ID); err != nil {
			return fiber.ErrInternalServerError
		}
		return ctx.SendStatus(204)
	}
}
