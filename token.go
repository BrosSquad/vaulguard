package main

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

func tokenCommands(ctx context.Context) *cobra.Command {
	t := &cobra.Command{
		Use: "token",
	}

	create := &cobra.Command{
		Use:  "create",
		Long: "Create new token for application",
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			app, err := applicationService.GetByName(ctx, args[0])

			if err != nil {
				log.Fatal(err)
			}

			tokenStr := tokenService.Generate(app.ID)

			if tokenStr == "" {
				log.Fatal("Error while generating Auth Token")
			}

			fmt.Printf("Auth Token: %s\n", tokenStr)
		},
	}

	t.AddCommand(create)

	return t
}
