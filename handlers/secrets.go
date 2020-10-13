package handlers

import (
	"github.com/BrosSquad/vaulguard/middleware"
	"github.com/BrosSquad/vaulguard/models"
	"github.com/BrosSquad/vaulguard/services/secret"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type secretHandlers struct {
	validator *validator.Validate
	service   secret.Service
}

func RegisterSecretHandlers(validate *validator.Validate, service secret.Service, r fiber.Router) {
	secretHandlers := secretHandlers{
		validator: validate,
		service:   service,
	}
	r.Get("/", secretHandlers.getSecrets).Use(middleware.ParsePageAndPerPage)
	r.Get("/many", secretHandlers.getManySecrets)
	r.Post("/", secretHandlers.createSecret)
	r.Delete("/invalidate", secretHandlers.invalidateCache)
}

func (s secretHandlers) getSecrets(ctx *fiber.Ctx) error {
	app := ctx.Locals("application").(models.ApplicationDto)
	page := ctx.Locals("page").(int)
	perPage := ctx.Locals("perPage").(int)

	secrets, err := s.service.Paginate(app.ID, page, perPage)

	if err != nil {
		return err
	}

	return ctx.JSON(secrets)
}

func (s secretHandlers) getManySecrets(c *fiber.Ctx) error {
	type query struct {
		keys []string `query:"keys"`
	}
	var keysStruct query
	app := c.Locals("application").(models.ApplicationDto)
	if err := c.QueryParser(&keysStruct); err != nil {
		return fiber.ErrBadRequest
	}

	secrets, err := s.service.Get(app.ID, keysStruct.keys)
	if err != nil {
		return err
	}

	return c.JSON(secrets)
}

func (s secretHandlers) createSecret(c *fiber.Ctx) error {
	type payload struct {
		Key   string `json:"key" validate:"required"`
		Value string `json:"value" validate:"required"`
	}

	var p payload
	app := c.Locals("application").(models.ApplicationDto)
	if err := c.BodyParser(&p); err != nil {
		return fiber.ErrBadRequest
	}

	if err := s.validator.Struct(p); err != nil {
		return err
	}

	data, err := s.service.Create(app.ID, p.Key, p.Value)

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(struct {
		ID    interface{} `json:"id"`
		Key   string      `json:"key"`
		Value string      `json:"value"`
	}{
		ID:    data.ID,
		Key:   data.Key,
		Value: p.Value,
	})
}

func (s secretHandlers) invalidateCache(ctx *fiber.Ctx) error {
	app := ctx.Locals("application").(models.ApplicationDto)

	if err := s.service.InvalidateCache(app.ID); err != nil {
		return fiber.ErrInternalServerError
	}
	return ctx.SendStatus(fiber.StatusNoContent)
}
