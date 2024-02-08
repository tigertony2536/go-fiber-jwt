package domain

import "gorm.io/gorm"

type (
	UserLogin struct {
		gorm.Model
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	Quote struct {
		Text []string `json:"text"`
	}

	JWT struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	Token struct {
		Email string
		Role  string
	}
)
