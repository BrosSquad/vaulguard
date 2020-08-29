package services

import (
	"os"
	"testing"

	"gorm.io/driver/sqlite"

	"github.com/BrosSquad/vaulguard/models"
	"gorm.io/gorm"
)

func TestToken(t *testing.T) {
	conn, err := gorm.Open(sqlite.Open("token_test.db"), &gorm.Config{})
	defer os.Remove("token_test.db")

	if err != nil {
		t.Fatal(err)
		return
	}

	if err := conn.AutoMigrate(&models.Application{}, &models.Token{}); err != nil {
		t.Fatal(err)
		return
	}

	app := models.Application{
		Name: "Test App",
	}
	conn.Create(&app)

	t.Run("Generate", func(t *testing.T) {
		s := NewTokenService(conn)
		_ = s.Generate(app.ID)
	})

	t.Run("Verify", func(t *testing.T) {
		s := NewTokenService(conn)

		token := s.Generate(app.ID)

		if !s.Verify(token) {
			t.Fatal("Token is not valid")
		}
	})
}
