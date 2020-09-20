package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/BrosSquad/vaulguard/middleware"
	"github.com/BrosSquad/vaulguard/services/application"
)

func RegisterApplicationHandlers(service application.Service, r fiber.Router) {
	r.Get("/", getApplications(service)).Use(middleware.ParsePageAndPerPage)
	r.Get("/search", searchApplications(service))
	r.Get("/:id", getApplication(service))
	r.Post("/", createApplication(service))
	r.Put("/:id", updateApplication(service))
	r.Delete("/:id", deleteApplication(service))
}

func getApplications(service application.Service) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		page := ctx.Locals("page").(int)
		perPage := ctx.Locals("perPage").(int)
		apps, err := service.Get(page, perPage)

		if err != nil {
			return fiber.ErrInternalServerError
		}

		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"data": apps,
		})
	}
}

func getApplication(service application.Service) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		return nil
	}
}

func searchApplications(service application.Service) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		return nil
	}
}

func createApplication(service application.Service) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		return nil
	}
}

func updateApplication(service application.Service) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		return nil
	}
}

func deleteApplication(service application.Service) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		return nil
	}
}
