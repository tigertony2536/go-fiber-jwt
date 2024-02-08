package domain

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func CreateToken(email, role string) (JWT, error) {
	var loginToken JWT
	t1 := jwt.New(jwt.SigningMethodHS256)
	c1 := t1.Claims.(jwt.MapClaims)
	c1["email"] = email
	c1["role"] = role
	c1["exp"] = time.Now().Add(time.Hour * 24).Unix()
	t, err := t1.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return loginToken, err
	}
	loginToken.AccessToken = t
	t2 := jwt.New(jwt.SigningMethodHS256)
	c2 := t2.Claims.(jwt.MapClaims)
	c2["email"] = email
	c2["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()
	rt, err := t2.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return loginToken, err
	}
	loginToken.RefreshToken = rt
	return loginToken, nil
}

func BinaryConvertor(number int, bits int) []int {
	result := make([]int, 0)
	for number > 0 {
		result = append(result, number%2)
		number /= 2
	}
	for i := len(result) - 1; len(result) != bits; i++ {
		result = append(result, 0)
	}
	return result
}
