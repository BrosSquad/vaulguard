package cmd

import (
	"context"
	"fmt"
	"github.com/BrosSquad/vaulguard/services/application"
	"github.com/BrosSquad/vaulguard/services/token"
	"github.com/spf13/cobra"
	"log"
)

type tokenCommand struct {
	ctx                context.Context
	applicationService application.Service
	tokenService       token.Service
}

func (tc tokenCommand) Execute(cmd *cobra.Command, args []string) error {
	app, err := tc.applicationService.GetByName(tc.ctx, args[0])

	if err != nil {
		return err
		//log.Fatal(err)
	}

	tokenStr := tc.tokenService.Generate(app.ID)

	if tokenStr == "" {
		log.Fatal("Error while generating Auth Token")
	}

	fmt.Printf("Auth Token: %s\n", tokenStr)

	return nil
}

func NewTokenCommand(ctx context.Context, applicationService application.Service, tokenService token.Service) Command {
	return tokenCommand{
		ctx:                ctx,
		applicationService: applicationService,
		tokenService:       tokenService,
	}
}
