package database_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tigertony2536/go-login/internal/adaptor/database"
	"github.com/tigertony2536/go-login/internal/config"
	"github.com/tigertony2536/go-login/internal/core/domain"
)

func TestCreateSession(t *testing.T) {
	cfg := config.NewConfig()
	gorm, db, err := database.NewTestDB()
	if err != nil {
		t.Fatal(err)
	}
	gorm.AutoMigrate(&domain.Session{})
	ar := database.NewAuthRepositoryImpl(gorm, cfg.HttpConfig)
	defer db.Close()
	t.Run("Create user successfully", func(t *testing.T) {
		expected := domain.Session{UserID: 1, RefreshToken: "token"}
		csErr := ar.CreateSession(&expected)

		result, err := ar.GetSessionByUserID(1)
		if err != nil {
			t.Fatal(err)
		}
		assert.NoError(t, csErr)
		assert.Equal(t, expected.UserID, result.UserID)
		assert.Equal(t, expected.RefreshToken, result.RefreshToken)
	})
}
