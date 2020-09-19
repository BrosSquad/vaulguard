package handlers

import (
	"errors"
	"github.com/BrosSquad/vaulguard/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Error(ctx *fiber.Ctx, err error) error {
	ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	if errors.Is(err, services.ErrAlreadyExists) {
		return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{"message": "Data already exists!"})
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Data not found!"})
	}

	// Send custom error page
	if err != nil {
		// In case the SendFile fails
		return ctx.Status(code).JSON(fiber.Map{"message": "An error has occurred!"})
	}

	// Return from handler
	return nil
}
