package services

import (
	"github.com/BrosSquad/vaulguard/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"testing"
)

func TestApplicationService(t *testing.T) {
	conn, err := gorm.Open(sqlite.Open("application_test.db"), &gorm.Config{})

	defer os.Remove("application_test.db")
	if err != nil {
		t.Error(err)
		return
	}

	if err := conn.AutoMigrate(&models.Application{}, &models.Token{}); err != nil {
		t.Error(err)
	}

	service := NewApplicationService(conn)

	t.Run("CreateApplication", func(t *testing.T) {
		app, err := service.Create("Test Application")

		if err != nil {
			t.Errorf("Error while inserting new application: %v", err)
		}

		if app.ID == 0 {
			t.Error("Application is not inserted")
		}
	})

	t.Run("Create2ApplicationWithSameName", func(t *testing.T) {
		if _, err := service.Create("Test Application"); err != nil {
			t.Errorf("Error while inserting new application: %v", err)
		}

		if _, err := service.Create("Test Application"); err != ErrAlreadyExists {
			t.Errorf("Inserting 2 applications with same name successed: %v", err)
		}
	})

	t.Run("DeleteApplication", func(t *testing.T) {
		app, err := service.Create("Test Application")

		if err != nil {
			t.Errorf("Error while inserting new application: %v", err)
		}

		if err := service.Delete(app.ID); err != nil {
			t.Errorf("Error while deleting application with ID: %d, %v", app.ID, err)
		}
	})

	t.Run("UpdateApplication", func(t *testing.T) {
		app, err := service.Create("Test Application")

		if err != nil {
			t.Errorf("Error while inserting new application: %v", err)
		}

		if app.Name != "Test Application" {
			t.Errorf("Application name is not one that is expected: %s", app.Name)
		}

		app, err = service.Update(app.ID, "Changed Name")

		if err != nil {
			t.Errorf("Error while updating application: %v", err)
		}

		if app.Name != "Changed Name" {
			t.Errorf("Application name is not one that is expected: %s", app.Name)
		}
	})

}
