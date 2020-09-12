package handlers

import (
	"strconv"

	"github.com/BrosSquad/vaulguard/models"
	"github.com/BrosSquad/vaulguard/services/secret"
	"github.com/gofiber/fiber"
)

func RegisterSecretHandlers(service secret.Service, r fiber.Router) {
	r.Get("/", getSecrets(service))
	r.Post("/many", getManySecrets(service))
	r.Delete("/invalidate", invalidateCache(service))
}

func getSecrets(service secret.Service) func(*fiber.Ctx) {
	return func(ctx *fiber.Ctx) {
		app := ctx.Locals("application").(models.ApplicationDto)
		pageStr := ctx.Query("page", "1")
		perPageStr := ctx.Query("perPage", "10")

		page, err := strconv.Atoi(pageStr)

		if err != nil {
			ctx.Next(fiber.NewError(400, "page query parameter is not a number"))
			return
		}

		perPage, err := strconv.Atoi(perPageStr)
		if err != nil {
			ctx.Next(fiber.NewError(400, "perPage query parameter is not a number"))
			return
		}

		secrets, err := service.Paginate(app.ID, page, perPage)

		if err != nil {
			ctx.Next(err)
			return
		}

		ctx.JSON(secrets)
	}
}

func getManySecrets(service secret.Service) func(*fiber.Ctx) {
	return func(ctx *fiber.Ctx) {
		var keys []string
		app := ctx.Locals("application").(models.ApplicationDto)
		if err := ctx.BodyParser(&keys); err != nil {
			ctx.Next(fiber.NewError(400, "Invalid Payload"))
			return
		}

		secrets, err := service.Get(app.ID, keys)

		if err != nil {
			ctx.Next(err)
			return
		}

		ctx.JSON(secrets)
	}
}

func invalidateCache(service secret.Service) func(*fiber.Ctx) {
	return func(ctx *fiber.Ctx) {
		app := ctx.Locals("application").(models.ApplicationDto)

		if err := service.InvalidateCache(app.ID); err != nil {
			ctx.Next(fiber.NewError(500, "Error while invalidating the cache"))
			return
		}

		ctx.SendStatus(204)
	}
}
