package middleware

import (
	"strings"

	"github.com/BrosSquad/vaulguard/services/token"
	"github.com/gofiber/fiber/v2"
)

type TokenAuthConfig struct {
	Header       string
	HeaderPrefix string
	TokenService token.Service
}

func TokenAuthMiddleware(config TokenAuthConfig) func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		authHeader := ctx.Get(config.Header)
		headerPrefixLen := len(config.HeaderPrefix)

		if authHeader == "" ||
			len(authHeader) < (headerPrefixLen+1) ||
			strings.ToLower(authHeader[0:headerPrefixLen]) != config.HeaderPrefix {
			return fiber.NewError(401, "Unauthorized")
		}

		t := authHeader[headerPrefixLen:]

		app, ok := config.TokenService.Verify(t)

		if !ok {
			return fiber.NewError(401, "Unauthorized")
		}

		ctx.Locals("application", app)
		return ctx.Next()
	}
}
