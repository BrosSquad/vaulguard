package middleware

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func ParsePageAndPerPage(ctx *fiber.Ctx) error {
	pageStr := ctx.Query("page", "1")
	perPageStr := ctx.Query("perPage", "10")

	page, err := strconv.ParseInt(pageStr, 10, 64)

	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "page query parameter is not a number")
	}

	perPage, err := strconv.ParseInt(perPageStr, 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "perPage query parameter is not a number")
	}

	ctx.Locals("page", int(page))
	ctx.Locals("perPage", int(perPage))

	return ctx.Next()
}
