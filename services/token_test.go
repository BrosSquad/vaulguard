package services

import (
	"os"
	"testing"

	"github.com/BrosSquad/vaulguard/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DbConn *gorm.DB
var app models.Application

func setupToken(t *testing.T) {
	var err error
	DbConn, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	if err != nil {
		teardownToken()
		t.Error(err)
	}

	if err := DbConn.AutoMigrate(&models.Application{}, &models.Token{}); err != nil {
		teardownToken()
		t.Error(err)
	}

	app = models.Application{
		Name: "Test App",
	}
	DbConn.Create(&app)
}

func teardownToken() {
	_ = os.Remove("test.db")
}

func TestTokenGenerate(t *testing.T) {
	setupToken(t)
	defer teardownToken()

	s := NewTokenService(DbConn)

	_ = s.Generate(app.ID)
}
