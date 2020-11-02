package token

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/BrosSquad/vaulguard/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestSqliteToken(t *testing.T) {
	ctx := context.Background()
	t.Parallel()
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
		s := NewService(NewSqlStorage(conn))
		_ = s.Generate(ctx, app.ID)
	})

	t.Run("Verify", func(t *testing.T) {
		s := NewService(NewSqlStorage(conn))
		token := s.Generate(ctx, app.ID)

		if _, ok := s.Verify(ctx, token); !ok {
			t.Fatal("Token is not valid")
		}
	})
}

func TestMongoToken(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	mongoURI := os.Getenv("VAULGUARD_MONGO_TESTING")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))

	if err != nil {
		t.Fatal(err)
	}
	if err := client.Connect(ctx); err != nil {
		t.Fatal(err)
	}
	db := client.Database("vaulguard_test")

	defer db.Drop(ctx)
	defer client.Disconnect(ctx)

	result, err := db.Collection("applications").InsertOne(ctx, bson.M{
		"Name":      "Test Application",
		"CreatedAt": time.Now(),
		"UpdatedAt": time.Now(),
	})

	if err != nil {
		t.Fatal(err)
	}

	service := NewService(NewMongoStorage(db.Collection("tokens")))

	t.Run("Generate", func(t *testing.T) {
		token := service.Generate(ctx, result.InsertedID)

		if token == "" {
			t.Fatal("Token is not generated")
		}

		values := strings.Split(token, ".")

		if len(values) != 3 && values[0] != "VaulGuard" {
			t.Fatal("Token is not valid")
		}
	})

	t.Run("Verify", func(t *testing.T) {
		token := service.Generate(ctx, result.InsertedID)
		if token == "" {
			t.Fatal("Token is not generated")
		}

		_, isValid := service.Verify(ctx, token)

		if !isValid {
			t.Fatal("Token should be valid")
		}
	})
}
