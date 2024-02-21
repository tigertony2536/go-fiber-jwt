package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

type (
	UserLogin struct {
		gorm.Model
		Email     string `gorm:"unique"`
		Password  string
		SessionID string
		Role      string
	}

	PairedToken struct {
		AccessToken  string
		RefreshToken string
	}

	JWTClaim struct {
		UserID uint
		Email  string
		Role   string
		jwt.RegisteredClaims
	}

	Session struct {
		SessionID    uint `gorm:"primarykey"`
		UserID       uint
		RefreshToken string
		CreatedAt    time.Time
	}
)
