package handlers

import (
	"strconv"

	"github.com/BrosSquad/vaulguard/models"
	"github.com/BrosSquad/vaulguard/services"
	"github.com/gofiber/fiber"
)

func RegisterSecretHandlers(service services.SecretService, r fiber.Router) {
	r.Get("/", getSecrets(service))
}

func getSecrets(service services.SecretService) func(*fiber.Ctx) {
	return func(ctx *fiber.Ctx) {
		app := ctx.Locals("application").(models.Application)
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

		secrets, err := service.Get(app.ID, page, perPage)

		if err != nil {
			ctx.Next(err)
			return
		}

		ctx.JSON(secrets)
	}
}