package core

import (
	"errors"
	"log"
	"runtime"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/tigertony2536/go-login/internal/core/domain"
)

func CreatePairedToken(user *domain.UserLogin, secret string) (*domain.PairedToken, error) {
	var token domain.PairedToken
	accessClaim := domain.JWTClaim{user.ID, user.Email, user.Role, jwt.RegisteredClaims{IssuedAt: jwt.NewNumericDate(time.Now())}}
	at := CreateToken(accessClaim, secret)
	token.AccessToken = at
	refreshClaim := domain.JWTClaim{user.ID, "", "", jwt.RegisteredClaims{IssuedAt: jwt.NewNumericDate(time.Now())}}
	rt := CreateToken(refreshClaim, secret)
	token.RefreshToken = rt
	return &token, nil
}

func CreateToken(claim domain.JWTClaim, secret string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	jwtToken, err := token.SignedString([]byte(secret))
	if err != nil {
		panic(err)
	}
	return jwtToken
}

func ParseToken(token, secret string) (*domain.JWTClaim, error) {
	var userClaim domain.JWTClaim
	t, err := jwt.ParseWithClaims(token, &userClaim, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !t.Valid {
		return nil, errors.New("invalid token")
	}
	return &userClaim, nil
}

func GetClaim(c *fiber.Ctx, secret string) *domain.JWTClaim {
	auth := c.Get("Authorization")
	token, _ := strings.CutPrefix(auth, "Bearer ")
	claim, err := ParseToken(token, secret)
	if err != nil {
		log.Fatal(err)
	}
	return claim
}

func GetFuncName(n int) string {
	pc, _, _, ok := runtime.Caller(n)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		name := details.Name()
		index := strings.LastIndex(name, ".")
		if index >= 0 {
			return name[index+1:]
		}
		return name
	}
	return ""
}
