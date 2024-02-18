package domain

import (
	"gorm.io/gorm"
)

type (
	UserLogin struct {
		gorm.Model
		Email        string `gorm:"unique"`
		Password     string
		RefreshToken string
		Role         string
	}

	Token struct {
		AccessToken  string
		RefreshToken string
	}
)
