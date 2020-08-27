package services

import (
	"os"
	"testing"

	"github.com/BrosSquad/vaulguard/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func teardownToken() {
	_ = os.Remove("test.db")
}

func setupToken(DbConn **gorm.DB, app *models.Application, t *testing.T) {
	var err error
	*DbConn, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	if err != nil {
		teardownToken()
		t.Error(err)
	}

	if err := (*DbConn).AutoMigrate(&models.Application{}, &models.Token{}); err != nil {
		teardownToken()
		t.Error(err)
	}

	*app = models.Application{
		Name: "Test App",
	}
	(*DbConn).Create(app)
}

func TestToken(t *testing.T) {
	var DbConn *gorm.DB
	var app models.Application
	setupToken(&DbConn, &app, t)
	defer teardownToken()

	t.Run("Generate", func(t *testing.T) {
		s := NewTokenService(DbConn)
		_ = s.Generate(app.ID)
	})

	t.Run("Verify", func(t *testing.T) {
		s := NewTokenService(DbConn)

		token := s.Generate(app.ID)

		if !s.Verify(token) {
			t.Error("Token is not valid")
		}
	})
}
