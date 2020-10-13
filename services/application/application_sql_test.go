package application

import (
	"errors"
	"github.com/BrosSquad/vaulguard/models"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"math/rand"
	"os"
	"testing"
)

func TestApplicationService(t *testing.T) {
	t.Parallel()
	asserts := require.New(t)
	conn, err := gorm.Open(sqlite.Open("application_test.db"), &gorm.Config{})
	db, err := conn.DB()
	asserts.Nil(err)
	defer os.Remove("application_test.db")
	defer db.Close()

	asserts.Nil(err)
	asserts.Nil(conn.AutoMigrate(&models.Application{}, &models.Token{}))
	service := NewSqlService(conn)

	t.Run("ListApplications", func(t *testing.T) {
		conn.Create(&models.Application{Name: "List App"})
		service.List(func(dtos []models.ApplicationDto) error {
			asserts.Len(dtos, 1)
			asserts.EqualValues("List App", dtos[0].Name)
			return nil
		})
	})

	t.Run("CreateApplication", func(t *testing.T) {
		app, err := service.Create("Test Application")
		asserts.Nil(err)
		asserts.Greater(app.ID, uint(0))
		asserts.EqualValues("Test Application", app.Name)
		asserts.False(app.CreatedAt.IsZero())
		asserts.False(app.UpdatedAt.IsZero())
	})

	t.Run("Create2ApplicationWithSameName", func(t *testing.T) {
		_, err := service.Create("Test Application 2")
		asserts.Nil(err)
		_, err = service.Create("Test Application 2")
		asserts.NotNil(err)
		//asserts.True(err == services.ErrAlreadyExists)
	})

	t.Run("DeleteApplication", func(t *testing.T) {
		app, err := service.Create("Test Application 3")
		asserts.Nil(err)
		asserts.Nil(service.Delete(app.ID))
	})

	t.Run("UpdateApplication", func(t *testing.T) {
		app, err := service.Create("Test Application 4")
		asserts.Nil(err)
		asserts.EqualValues("Test Application 4", app.Name)
		app, err = service.Update(app.ID, "Changed Name")
		asserts.Nil(err)
		asserts.EqualValues("Changed Name", app.Name)
	})

	t.Run("UpdateApplicationNotFound", func(t *testing.T) {
		_, err := service.Update(uint(152000), "New Name")
		asserts.NotNil(err)
		asserts.True(errors.Is(err, gorm.ErrRecordNotFound))
	})

	t.Run("GetByName", func(t *testing.T) {
		appNames := []string{"Test App1", "Test App2", "Test App3"}
		for _, appName := range appNames {
			_, err := service.Create(appName)
			asserts.Nil(err)
		}

		app, err := service.GetByName("Test App2")
		asserts.Nil(err)
		asserts.Greater(app.ID, uint(0))
		asserts.EqualValues("Test App2", app.Name)
		asserts.False(app.CreatedAt.IsZero())
		asserts.False(app.UpdatedAt.IsZero())
	})

	t.Run("GetOne", func(t *testing.T) {
		var app models.ApplicationDto
		appNames := []string{"Test GetOne App1", "Test GetOne App2", "Test GetOne App3"}
		randomInt := rand.Int31n(int32(len(appNames)))
		for i, appName := range appNames {
			a, err := service.Create(appName)
			asserts.Nil(err)

			if randomInt == int32(i) {
				app = a
			}
		}

		a, err := service.GetOne(app.ID)
		asserts.Nil(err)
		asserts.Greater(a.ID, uint(0))
		asserts.EqualValues(app.Name, a.Name)
		asserts.False(a.CreatedAt.IsZero())
		asserts.False(a.UpdatedAt.IsZero())
	})

	t.Run("GetOneNotFound", func(t *testing.T) {
		_, err := service.GetOne(uint(152000))
		asserts.NotNil(err)
		asserts.True(errors.Is(err, gorm.ErrRecordNotFound))
	})

	t.Run("Get", func(t *testing.T) {
		asserts.Nil(conn.Delete(&models.Application{}, " 1 = 1").Error)
		appNames := []string{"Test Get App1", "Test Get App2", "Test Get App3", "Test Get App4"}
		for _, appName := range appNames {
			_, err := service.Create(appName)
			asserts.Nil(err)
		}

		apps, err := service.Get(1, 3)
		asserts.Nil(err)
		asserts.Len(apps, 3)

		for _, app := range apps {
			asserts.Greater(app.ID, uint(0))
			asserts.Contains(appNames, app.Name)
			asserts.False(app.CreatedAt.IsZero())
			asserts.False(app.UpdatedAt.IsZero())
		}
	})

	t.Run("GetSecondPage", func(t *testing.T) {
		asserts.Nil(conn.Delete(&models.Application{}, " 1 = 1").Error)
		appNames := []string{"Test Get App1", "Test Get App2", "Test Get App3", "Test Get App4"}
		for _, appName := range appNames {
			_, err := service.Create(appName)
			asserts.Nil(err)
		}

		apps, err := service.Get(2, 3)
		asserts.Nil(err)
		asserts.Len(apps, 1)

		for _, app := range apps {
			asserts.Greater(app.ID, uint(0))
			asserts.Contains(appNames, app.Name)
			asserts.False(app.CreatedAt.IsZero())
			asserts.False(app.UpdatedAt.IsZero())
		}
	})

	t.Run("GetWithNegativePageAndPerPage", func(t *testing.T) {
		asserts.Nil(conn.Delete(&models.Application{}, " 1 = 1").Error)
		appNames := []string{"Test Get App1", "Test Get App2", "Test Get App3", "Test Get App4"}
		for _, appName := range appNames {
			_, err := service.Create(appName)
			asserts.Nil(err)
		}

		apps, err := service.Get(-1, -3)
		asserts.Nil(err)
		asserts.Len(apps, 3)

		for _, app := range apps {
			asserts.Greater(app.ID, uint(0))
			asserts.Contains(appNames, app.Name)
			asserts.False(app.CreatedAt.IsZero())
			asserts.False(app.UpdatedAt.IsZero())
		}
	})
}
