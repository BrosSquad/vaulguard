package services

import (
	"os"
	"testing"

	"github.com/BrosSquad/vaulguard/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestToken(t *testing.T) {
	conn, err := gorm.Open(sqlite.Open("token_test.db"), &gorm.Config{})
	db, _ := conn.DB()
	defer os.Remove("token_test.db")
	defer db.Close()

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

		if _, ok := s.Verify(token); !ok {
			t.Fatal("Token is not valid")
		}
	})
}
