package handlers

import (
	"strconv"

	"github.com/BrosSquad/vaulguard/models"
	"github.com/BrosSquad/vaulguard/services/secret"
	"github.com/gofiber/fiber/v2"
)

func RegisterSecretHandlers(service secret.Service, r fiber.Router) {
	r.Get("/", getSecrets(service))
	r.Get("/many", getManySecrets(service))
	r.Delete("/invalidate", invalidateCache(service))
}

func getSecrets(service secret.Service) func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		app := ctx.Locals("application").(models.ApplicationDto)
		pageStr := ctx.Query("page", "1")
		perPageStr := ctx.Query("perPage", "10")

		page, err := strconv.Atoi(pageStr)

		if err != nil {
			return fiber.NewError(400, "page query parameter is not a number")
		}

		perPage, err := strconv.Atoi(perPageStr)
		if err != nil {
			return fiber.NewError(400, "perPage query parameter is not a number")
		}

		secrets, err := service.Paginate(app.ID, page, perPage)

		if err != nil {
			return err
		}

		return ctx.JSON(secrets)
	}
}

func getManySecrets(service secret.Service) func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		keysStruct := struct {
			keys []string `query:"keys"`
		}{}
		app := ctx.Locals("application").(models.ApplicationDto)
		if err := ctx.QueryParser(&keysStruct); err != nil {
			return fiber.NewError(400, "Invalid Payload")
		}

		secrets, err := service.Get(app.ID, keysStruct.keys)
		if err != nil {
			return err
		}

		return ctx.JSON(secrets)
	}
}

func invalidateCache(service secret.Service) func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		app := ctx.Locals("application").(models.ApplicationDto)

		if err := service.InvalidateCache(app.ID); err != nil {
			return fiber.NewError(500, "Error while invalidating the cache")
		}
		return ctx.SendStatus(204)
	}
}
