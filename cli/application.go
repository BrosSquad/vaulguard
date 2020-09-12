package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/BrosSquad/vaulguard/models"
	"github.com/spf13/cobra"
)

func applicationCommands() *cobra.Command {
	app := &cobra.Command{
		Use: "app",
	}

	list := &cobra.Command{
		Use:  "list",
		Long: "List all applications in the database",
		Run: func(cmd *cobra.Command, args []string) {
			iterate := func(apps []models.ApplicationDto) error {
				for _, app := range apps {
					fmt.Printf("ID: %d, Name: %s\n", app.ID, app.Name)
				}

				return nil
			}
			err := applicationService.List(iterate)

			if err != nil {
				log.Fatal(err)
			}
		},
	}

	findByName := &cobra.Command{
		Use:  "by-name",
		Long: "Search for application by it's name",
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			app, err := applicationService.GetByName(args[0])

			if err != nil {
				log.Fatal(err.Error())
			}

			fmt.Printf("ID: %d, Name:%s\n", app.ID, app.Name)
		},
	}

	create := &cobra.Command{
		Use:  "create",
		Long: "Create new application for VaulGuard",
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			app, err := applicationService.Create(args[0])

			if err != nil {
				log.Fatal(err.Error())
			}

			tokenStr := tokenService.Generate(app.ID)

			if tokenStr == "" {
				log.Fatal("Error while generating Auth Token")
			}

			fmt.Printf("New Application created: ID: %d Name: %s\n", app.ID, app.Name)
			fmt.Printf("Auth Token: %s\n", tokenStr)
		},
	}

	del := &cobra.Command{
		Use:  "delete",
		Long: "Delete application with given application ID",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("Not enough argument passed to delete command")
			}
			_, err := strconv.Atoi(args[1])

			if err != nil {
				return errors.New("Argument is not of type number")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			id, _ := strconv.Atoi(args[1])
			err := applicationService.Delete(uint(id))

			if err != nil {
				log.Fatal(err.Error())
			}

			fmt.Println("Application successfully deleted")
		},
	}

	app.AddCommand(list)
	app.AddCommand(findByName)
	app.AddCommand(create)
	app.AddCommand(del)

	return app
}
