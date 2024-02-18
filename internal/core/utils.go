package core

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/tigertony2536/go-login/internal/core/domain"
)

func CreateToken(email, role string) (domain.Token, error) {
	var token domain.Token
	t1 := jwt.New(jwt.SigningMethodHS256)
	c1 := t1.Claims.(jwt.MapClaims)
	c1["email"] = email
	c1["role"] = role
	c1["exp"] = time.Now().Add(time.Hour * 24).Unix()
	t, err := t1.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return token, err
	}
	token.AccessToken = t
	t2 := jwt.New(jwt.SigningMethodHS256)
	c2 := t2.Claims.(jwt.MapClaims)
	c2["email"] = email
	c2["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()
	rt, err := t2.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return token, err
	}
	token.RefreshToken = rt
	return token, nil
}
