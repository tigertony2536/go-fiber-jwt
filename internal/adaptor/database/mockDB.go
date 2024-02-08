package database

import (
	"errors"

	"github.com/tigertony2536/go-login/internal/domain"
)

type MockDB struct{}

func NewMockDB() *MockDB {
	return &MockDB{}
}

var users = []domain.UserLogin{
	{Email: "user1@user.com", Password: "user1"},
	{Email: "user2@user.com", Password: "user2"},
}

func (*MockDB) GetUserByEmil(email string) (domain.UserLogin, error) {
	var user domain.UserLogin
	for _, u := range users {
		if email == u.Email {
			return u, nil
		}
	}
	return user, errors.New("user not found")
}
