package services

import (
	"crypto/rand"
	"os"
	"testing"

	"github.com/BrosSquad/vaulguard/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestNewGormSecretStorage(t *testing.T) {
	conn, err := gorm.Open(sqlite.Open("secret_test.db"), &gorm.Config{})

	defer os.Remove("secret_test.db")
	if err != nil {
		t.Error(err)
		return
	}

	if err := conn.AutoMigrate(&models.Application{}, &models.Token{}, &models.Secret{}); err != nil {
		t.Error(err)
	}

	application := models.Application{Name: "Test Application"}
	conn.Create(&application)

	key := make([]byte, 32)
	_, _ = rand.Read(key)
	encryptionService, _ := NewEncryptionService(key)
	service := NewGormSecretStorage(conn, encryptionService)

	t.Run("CreateSecret", func(t *testing.T) {
		value := "mysql://localhost:3306/database"
		_, err := service.Create(application.ID, "DATABASE_CONNECTION", value)

		if err != nil {
			t.Fatalf("Error while inserting new secret: %v", err)
		}

		secret, err := service.GetOne(application.ID, "DATABASE_CONNECTION")

		if err != nil {
			t.Fatalf("Error while geting secret: %v", err)
		}

		if secret.Value != value {
			t.Fatalf("Expected: %s, GOT: %s", value, secret.Value)
		}
	})

	t.Run("UpdateSecret", func(t *testing.T) {
		value := "mysql://localhost:3306/database"
		_, err := service.Create(application.ID, "DATABASE_CONNECTION_2", value)
		if err != nil {
			t.Fatalf("Error while inserting new secret: %v", err)
		}

		newValue := "postgres://localhost:5432/database"
		_, err = service.Update(application.ID, "DATABASE_CONNECTION_2", newValue)

		if err != nil {
			t.Fatalf("Error while updating secret: %v", err)
		}

		secretDecrypted, err := service.GetOne(application.ID, "DATABASE_CONNECTION_2")

		if err != nil {
			t.Fatalf("Secret with name `DATABASE_CONNECTION_2` does not exist: %v", err)
		}

		if secretDecrypted.Value != newValue {
			t.Fatal("Updating secert failed")
		}

		if secretDecrypted.Value == value {
			t.Fatal("Secret remained the same value as before")
		}
	})
}
