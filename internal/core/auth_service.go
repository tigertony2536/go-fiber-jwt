package core

import (
	"errors"

	"github.com/tigertony2536/go-login/internal/adaptor/database"
	"github.com/tigertony2536/go-login/internal/core/domain"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthServiceImpl struct {
	repo *database.UserRepositoryImpl
}

func NewAuthService(repo *database.UserRepositoryImpl) *AuthServiceImpl {
	return &AuthServiceImpl{repo: repo}
}

func (us *AuthServiceImpl) Register(user *domain.UserLogin) error {
	result, err := us.repo.GetUserByEmail(user.Email)
	var nilpointer *domain.UserLogin
	if errors.Is(err, gorm.ErrRecordNotFound) && result == nilpointer {
		user.Role = "user"
		_, err = us.repo.CreateUser(user)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("Register failed. This email have been used")
}

func (us *AuthServiceImpl) Login(req *domain.UserLogin) (*domain.Token, error) {
	u, err := us.repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 8)
	if err != nil {
		return nil, err
	}
	if bcrypt.CompareHashAndPassword(hash, []byte(u.Password)) != nil {
		return nil, err
	}

	token, err := CreateToken(req.Email, "users")
	if req.Password != u.Password {
		if err != nil {
			return nil, errors.New("create token failed")
		}
	}
	return &token, nil
}
