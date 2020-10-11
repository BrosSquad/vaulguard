package main

import (
	"github.com/BrosSquad/vaulguard/services/application"
	"github.com/BrosSquad/vaulguard/services/secret"
	"github.com/BrosSquad/vaulguard/services/token"
	"github.com/spf13/cobra"
)

var (
	applicationService application.Service
	secretService      secret.Service
	tokenService       token.Service
)

var (
	rootCmd *cobra.Command
)

func parseCommands() error {
	rootCmd = &cobra.Command{
		Use:   "vaulguard",
		Short: "VaulGuard CLI",
		Long:  "Command line interface for VaulGuard secret storage",
	}

	rootCmd.AddCommand(applicationCommands())
	rootCmd.AddCommand(tokenCommands())

	return rootCmd.Execute()
}

func main() {
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
