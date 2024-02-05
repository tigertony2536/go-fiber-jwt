package main

type (
	Quote struct {
		Text []string `json:"text"`
	}

	RequestLogin struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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
