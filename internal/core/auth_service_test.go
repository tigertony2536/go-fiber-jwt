package core_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tigertony2536/go-login/internal/adaptor/database"
	"github.com/tigertony2536/go-login/internal/core"
	"github.com/tigertony2536/go-login/internal/core/domain"
)

func TestRegister(t *testing.T) {
	gorm, db := database.InnitializeDB()
	ur := database.NewUserRepositoryImpl(gorm)
	as := core.NewAuthService(ur)
	defer db.Close()

	t.Run("Registered successfully", func(t *testing.T) {
		regisUser := domain.UserLogin{Email: "user1@user.com", Password: "user1"}
		var nilpointer *domain.UserLogin
		result, err := ur.GetUserByEmail(regisUser.Email)
		if err != nil && result == nilpointer {
			err := as.Register(&regisUser)
			assert.NoError(t, err)

			result, err := ur.GetUserByEmail("user1@user.com")
			if err != nil {
				t.Fatal("Registered user not found")
			}
			assert.Equal(t, regisUser.Email, result.Email)
			assert.Equal(t, regisUser.Password, result.Password)
		}
	})
}
