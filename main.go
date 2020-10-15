package main

import (
	"context"
	"github.com/BrosSquad/vaulguard/cmd"
	"github.com/BrosSquad/vaulguard/services/application"
	"github.com/BrosSquad/vaulguard/services/secret"
	"github.com/BrosSquad/vaulguard/services/token"
	"github.com/spf13/cobra"
	"log"
)

var (
	applicationService application.Service
	secretService      secret.Service
	tokenService       token.Service
)

var (
	rootCmd *cobra.Command
)

func createTokenCommand(tc cmd.Command, command *cobra.Command) *cobra.Command {
	create := &cobra.Command{
		Use:  "create",
		Long: "Create new token for application",
		Args: cobra.MinimumNArgs(1),
		RunE: tc.Execute,
	}
	command.AddCommand(create)

	return command
}

func main() {
	ctx := context.Background()
	rootCmd = &cobra.Command{
		Use:   "vaulguard",
		Short: "VaulGuard CLI",
		Long:  "Command line interface for VaulGuard secret storage",
	}

	rootCmd.AddCommand(createTokenCommand(
		cmd.NewTokenCommand(ctx, applicationService, tokenService),
		&cobra.Command{
			Use: "token",
		},
	))

	rootCmd.AddCommand(applicationCommands(ctx))
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Command error: %v", err)
	}
	//cfg, err := config.NewConfig(true)
	//
	//if err != nil {
	//	log.Fatalf("Error while creating app configuration: %v\n", err)
	//}
	//
	//conn, err := db.ConnectToDatabaseProvider(cfg.Database, cfg.DatabaseDSN)
	//
	//if err != nil {
	//	log.Fatalf("Error while connection to PostgreSQL: %v", err)
	//}
	//
	//if err := db.Migrate(); err != nil {
	//	log.Fatalf("Auto migration failed: %v", err)
	//}
	//
	//applicationService = application.NewSqlService(conn)
	//tokenService = token.NewSqlService(conn)
	//
	//if err := parseCommands(); err != nil {
	//	log.Fatal(err)
	//}
}
