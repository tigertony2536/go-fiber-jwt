package core_test

import (
	"log"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/tigertony2536/go-login/internal/adaptor/database"
	"github.com/tigertony2536/go-login/internal/adaptor/router"
	"github.com/tigertony2536/go-login/internal/config"
	"github.com/tigertony2536/go-login/internal/core"
	"github.com/tigertony2536/go-login/internal/core/domain"
)

func init() {
	err := godotenv.Load("D:\\dev\\go\\src\\03-side-projects\\go-login\\.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func TestRegister(t *testing.T) {
	cfg := config.NewConfig()
	gorm, db, err := database.NewTestDB()
	if err != nil {
		t.Fatal(err)
	}
	gorm.AutoMigrate(&domain.Session{})
	defer db.Close()
	ur := database.NewUserRepositoryImpl(gorm)
	ar := database.NewAuthRepositoryImpl(gorm, cfg.HttpConfig)
	as := core.NewAuthService(ur, ar)
	defer db.Close()

	t.Run("Registered successfully", func(t *testing.T) {
		regisUser := domain.UserLogin{Email: "user1@user.com", Password: "user1"}
		var nilpointer *domain.UserLogin
		result, err := ur.GetUserByEmail(regisUser.Email)
		if err != nil && result == nilpointer {
			err := as.Register(&regisUser, "user")
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

func TestIsExpired(t *testing.T) {
	cfg := config.NewConfig()
	gorm, db, err := database.NewTestDB()
	if err != nil {
		t.Fatal(err)
	}
	gorm.AutoMigrate(&domain.Session{})
	defer db.Close()

	ur := database.NewUserRepositoryImpl(gorm)
	ar := database.NewAuthRepositoryImpl(gorm, cfg.HttpConfig)
	as := core.NewAuthService(ur, ar)
	a := router.NewAuthHandler(as, ar)

	session := domain.Session{UserID: 1, RefreshToken: "Refresh", CreatedAt: time.Now()}
	err = a.Ar.CreateSession(&session)
	if err != nil {
		t.Fatal("Can not Create Session")
	}
	t.Run("Token is not expired", func(t *testing.T) {
		issuedAt := jwt.NewNumericDate(session.CreatedAt)
		claim := domain.JWTClaim{1, "user1@user.com", "user", jwt.RegisteredClaims{IssuedAt: issuedAt}}
		expired, err := a.As.IsExpired(claim.IssuedAt.Time)
		assert.NoError(t, err)
		assert.False(t, expired)
	})
	t.Run("Token is expired", func(t *testing.T) {
		issuedAt := jwt.NewNumericDate(session.CreatedAt.Add(-1 * 999 * time.Hour))
		claim := domain.JWTClaim{1, "user1@user.com", "user", jwt.RegisteredClaims{IssuedAt: issuedAt}}
		expired, err := a.As.IsExpired(claim.IssuedAt.Time)
		assert.NoError(t, err)
		assert.True(t, expired)
	})
}
