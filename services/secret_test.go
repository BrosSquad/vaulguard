package services

import (
	"crypto/rand"
	"encoding/base64"
	mathrand "math/rand"
	"os"
	"testing"

	"github.com/BrosSquad/vaulguard/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestNewGormSecretStorage(t *testing.T) {
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
			service.Create(application.ID, key, value)
		}

		secrets, err := service.Get(application.ID, []string{"SECRET_1", "SECRET_2", "SECRET_6"})

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
		_, err := service.Create(application.ID, "DATABASE_CONNECTION_2", value)
		if err != nil {
			t.Fatalf("Error while inserting new secret: %v", err)
		}

		newValue := "postgres://localhost:5432/database"
		_, err = service.Update(application.ID, "DATABASE_CONNECTION_2", "DATABASE_CONNECTION_2", newValue)

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

func BenchmarkSecretsInSqlite(b *testing.B) {
	conn, err := gorm.Open(postgres.New(postgres.Config{
		DSN: "host=localhost user=postgres password=postgres dbname=vaulguard_benchmark port=5432 sslmode=disable TimeZone=UTC",
	}), &gorm.Config{})
	db, _ := conn.DB()

	defer db.Close()
	if err != nil {
		b.Fatal(err)
		return
	}

	defer func() {
		conn.Exec("DELETE FROM secrets")
		conn.Exec("DELETE FROM applications")
	}()

	if err := conn.AutoMigrate(&models.Application{}, &models.Token{}, &models.Secret{}); err != nil {
		b.Fatal(err)
	}

	application := models.Application{Name: "Test Application"}
	conn.Create(&application)

	appKey := make([]byte, 32)
	_, _ = rand.Read(appKey)
	encryptionService, _ := NewEncryptionService(appKey)
	service := NewGormSecretStorage(conn, encryptionService)

	secretsMap := make([]Secret, 10000)

	keyBytes := make([]byte, 64)
	valueBytes := make([]byte, 64)

	for i := 0; i < 10000; i++ {
		_, _ = rand.Read(keyBytes)
		_, _ = rand.Read(valueBytes)
		key := base64.RawStdEncoding.EncodeToString(keyBytes)
		value := base64.RawStdEncoding.EncodeToString(valueBytes)

		secretsMap[i] = Secret{
			key, value,
		}
	}

	b.Run("Insert 10_000", func(b *testing.B) {
		for _, value := range secretsMap {
			service.Create(application.ID, value.Key, value.Value)
		}
	})

	b.Run("GetMultipleValues", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			length := mathrand.Int31n(int32(b.N))
			keys := make([]string, length)
			for j := 0; j < int(length); j++ {
				keys[j] = secretsMap[j].Key
			}

			_, err = service.Get(application.ID, keys)

			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
