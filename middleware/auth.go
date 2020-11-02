package middleware

import (
	"strings"

	"github.com/BrosSquad/vaulguard/services/token"
	"github.com/gofiber/fiber/v2"
)

type TokenAuthConfig struct {
	Headers        []string
	HeaderPrefixes []string
	TokenServices  []token.Service
}

func extractToken(header, prefix string, prefixLen int) (string, error) {
	if header == "" ||
		len(header) < (prefixLen+1) ||
		strings.ToLower(header[0:prefixLen]) != prefix {
		return "", fiber.ErrUnauthorized
	}

	return header[prefixLen:], nil
}

func TokenAuth(config TokenAuthConfig) fiber.Handler {
	headersLen := len(config.Headers)
	headerPrefixesLen := len(config.HeaderPrefixes)
	servicesLen := len(config.TokenServices)

	if headersLen != headerPrefixesLen || headersLen != servicesLen || headerPrefixesLen != servicesLen {
		panic("config.Headers, config.HeaderPrefixes and config.TokenServices must have same length")
	}

	return func(ctx *fiber.Ctx) error {
		for i, service := range config.TokenServices {
			t, err := extractToken(ctx.Get(config.Headers[i]), config.HeaderPrefixes[i], len(config.HeaderPrefixes[i]))

			if err != nil {
				return err
			}
			app, ok := service.Verify(ctx.Context(), t)
			if ok {
				ctx.Locals("application", app)
				return ctx.Next()
			}
		}

		return fiber.ErrUnauthorized

	}
}
