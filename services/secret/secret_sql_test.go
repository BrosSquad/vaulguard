package secret

import (
	"context"
	"crypto/rand"
	"os"
	"testing"

	"github.com/BrosSquad/vaulguard/models"
	"github.com/BrosSquad/vaulguard/services"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestNewGormSecretStorage(t *testing.T) {
	ctx := context.Background()
	conn, err := gorm.Open(sqlite.Open("secret_test.db"), &gorm.Config{})
	db, _ := conn.DB()

	defer os.Remove("secret_test.db")
	defer db.Close()
	if err != nil {
		t.Fatal(err)
		return
	}

	if err := conn.AutoMigrate(&models.Application{}, &models.Token{}, &models.Secret{}); err != nil {
		t.Fatal(err)
	}

	application := models.Application{Name: "Test Application"}
	conn.Create(&application)

	key := make([]byte, 32)
	_, _ = rand.Read(key)
	encryptionService, _ := services.NewSecretKeyEncryption(key)
	service := NewGormSecretStorage(GormSecretConfig{
		Encryption: encryptionService,
		DB:         conn,
		CacheSize:  32,
	})

	t.Run("CreateSecret", func(t *testing.T) {
		value := "mysql://localhost:3306/database"
		_, err := service.Create(ctx, application.ID, "DATABASE_CONNECTION", value)

		if err != nil {
			t.Fatalf("Error while inserting new secret: %v", err)
		}

		secret, err := service.GetOne(ctx, application.ID, "DATABASE_CONNECTION")

		if err != nil {
			t.Fatalf("Error while geting secret: %v", err)
		}

		if secret.Value != value {
			t.Fatalf("Expected: %s, GOT: %s", value, secret.Value)
		}
	})

	t.Run("MultipleSecretsReturned", func(t *testing.T) {
		secretsMap := map[string]string{
			"SECRET_1": "TEST",
			"SECRET_2": "TEST",
			"SECRET_3": "TEST",
			"SECRET_4": "TEST",
			"SECRET_5": "TEST",
			"SECRET_6": "TEST",
			"SECRET_7": "TEST",
		}

		for key, value := range secretsMap {
			_, err := service.Create(ctx, application.ID, key, value)
			if err != nil {
				t.Fatal(err)
			}
		}

		secrets, err := service.Get(ctx, application.ID, []string{"SECRET_1", "SECRET_2", "SECRET_6"})

		if err != nil {
			t.Fatalf("Error while fetching secrets: %v\n", err)
		}

		for key, secret := range secrets {
			if _, ok := secretsMap[key]; !ok {
				t.Fatalf("Secret %s does not exist\n", key)
			}

			if secret != secretsMap[key] {
				t.Fatalf("Secret %s does not have the same value as in MAP: %s\n", key, secret)
			}
		}

	})

	t.Run("UpdateSecret", func(t *testing.T) {
		value := "mysql://localhost:3306/database"
		_, err := service.Create(ctx, application.ID, "DATABASE_CONNECTION_2", value)
		if err != nil {
			t.Fatalf("Error while inserting new secret: %v", err)
		}

		newValue := "postgres://localhost:5432/database"
		_, err = service.Update(ctx, application.ID, "DATABASE_CONNECTION_2", "DATABASE_CONNECTION_2", newValue)

		if err != nil {
			t.Fatalf("Error while updating secret: %v", err)
		}

		secretDecrypted, err := service.GetOne(ctx, application.ID, "DATABASE_CONNECTION_2")

		if err != nil {
			t.Fatalf("Secret with name `DATABASE_CONNECTION_2` does not exist: %v", err)
		}

		if secretDecrypted.Value != newValue {
			t.Fatal("Updating secret failed")
		}

		if secretDecrypted.Value == value {
			t.Fatal("Secret remained the same value as before")
		}
	})
}
