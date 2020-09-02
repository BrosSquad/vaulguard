package services

import (
	"os"
	"testing"

	"github.com/BrosSquad/vaulguard/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestApplicationService(t *testing.T) {
	conn, err := gorm.Open(sqlite.Open("application_test.db"), &gorm.Config{})
	db, _ := conn.DB()

	defer os.Remove("application_test.db")
	defer db.Close()

	if err != nil {
		t.Fatal(err)
		return
	}

	if err := conn.AutoMigrate(&models.Application{}, &models.Token{}); err != nil {
		t.Fatal(err)
	}

	service := NewApplicationService(conn)

	t.Run("CreateApplication", func(t *testing.T) {
		app, err := service.Create("Test Application")

		if err != nil {
			t.Fatalf("Error while inserting new application: %v", err)
		}

		if app.ID == 0 {
			t.Fatal("Application is not inserted")
		}

		if app.Name != "Test Application" {
			t.Fatal("Application name not expected")
		}
	})

	t.Run("Create2ApplicationWithSameName", func(t *testing.T) {
		if _, err := service.Create("Test Application 2"); err != nil {
			t.Fatalf("Error while inserting new application: %v", err)
		}

		if _, err := service.Create("Test Application"); err != ErrAlreadyExists {
			t.Fatalf("Inserting 2 applications with same name successed: %v", err)
		}
	})

	t.Run("DeleteApplication", func(t *testing.T) {
		app, err := service.Create("Test Application 3")

		if err != nil {
			t.Fatalf("Error while inserting new application: %v", err)
		}

		if err := service.Delete(app.ID); err != nil {
			t.Fatalf("Error while deleting application with ID: %d, %v", app.ID, err)
		}
	})

	t.Run("UpdateApplication", func(t *testing.T) {
		app, err := service.Create("Test Application 4")

		if err != nil {
			t.Fatalf("Error while inserting new application: %v", err)
		}

		if app.Name != "Test Application 4" {
			t.Fatalf("Application name is not one that is expected: %s", app.Name)
		}

		app, err = service.Update(app.ID, "Changed Name")

		if err != nil {
			t.Fatalf("Error while updating application: %v", err)
		}

		if app.Name != "Changed Name" {
			t.Fatalf("Application name is not one that is expected: %s", app.Name)
		}
	})

}
