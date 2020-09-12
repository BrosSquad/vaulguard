package middleware

import (
	"strings"

	"github.com/BrosSquad/vaulguard/services/token"
	"github.com/gofiber/fiber"
)

type TokenAuthConfig struct {
	Header       string
	HeaderPrefix string
	TokenService token.Service
}

func TokenAuthMiddleware(config TokenAuthConfig) func(*fiber.Ctx) {
	return func(ctx *fiber.Ctx) {
		authHeader := ctx.Get(config.Header)
		headerPrefixLen := len(config.HeaderPrefix)
		if authHeader == "" ||
			len(authHeader) < (headerPrefixLen+1) ||
			strings.ToLower(authHeader[0:headerPrefixLen]) != config.HeaderPrefix {
			ctx.Next(fiber.NewError(401, "Unauthorized"))
			return
		}

		t := authHeader[headerPrefixLen:]

		app, ok := config.TokenService.Verify(t)

		if !ok {
			ctx.Next(fiber.NewError(401, "Unauthorized"))
			return
		}

		ctx.Locals("application", app)
		ctx.Next()
	}
}
