package api

import (
	"context"

	"github.com/BrosSquad/vaulguard/handlers"
	"github.com/BrosSquad/vaulguard/log"
	"github.com/BrosSquad/vaulguard/services/application"
	"github.com/BrosSquad/vaulguard/services/secret"
	"github.com/BrosSquad/vaulguard/services/token"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/BrosSquad/vaulguard/config"
	"github.com/BrosSquad/vaulguard/middleware"
)

type Interface interface {
	RegisterHandlers()
}

type Fiber struct {
	Ctx                   context.Context
	Cfg                   *config.Config
	App                   fiber.Router
	TokenCollection       *mongo.Collection
	SecretCollection      *mongo.Collection
	ApplicationCollection *mongo.Collection

	TokenService       token.Service
	ApplicationService application.Service
	SecretService      secret.Service
	Logger             *log.Logger
	Validator          *validator.Validate
}

func (f Fiber) RegisterHandlers() {
	f.registerSecrets()
	f.registerApplications()
}

func (f Fiber) registerApplications() {
	f.Logger.Debug("Starting to add APPLICATION routes.")
	handlers.RegisterApplicationHandlers(f.ApplicationService, f.App.Group("/applications"))
	f.Logger.Debug("APPLICATION routes added.")
}

func (f Fiber) registerSecrets() {
	f.Logger.Debug("Starting to add SECRET routes.")
	secretsGroup := f.App.Group("/secrets")

	secretsGroup.Use(middleware.TokenAuth(middleware.TokenAuthConfig{
		TokenServices:  []token.Service{f.TokenService},
		Headers:        []string{"authorization"},
		HeaderPrefixes: []string{"token "},
	}))
	handlers.RegisterSecretHandlers(f.Validator, f.SecretService, secretsGroup)

	f.Logger.Debug("SECRET routes added.")

}
