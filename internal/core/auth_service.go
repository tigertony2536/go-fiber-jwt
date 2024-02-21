package core

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/tigertony2536/go-login/internal/adaptor/database"
	"github.com/tigertony2536/go-login/internal/core/domain"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthServiceImpl struct {
	Repo *database.UserRepositoryImpl
	Ar   *database.AuthRepositoryImpl
}

func NewAuthService(repo *database.UserRepositoryImpl, ar *database.AuthRepositoryImpl) *AuthServiceImpl {
	return &AuthServiceImpl{Repo: repo, Ar: ar}
}

func (us *AuthServiceImpl) Register(user *domain.UserLogin, role string) error {
	result, err := us.Repo.GetUserByEmail(user.Email)
	var nilpointer *domain.UserLogin
	if errors.Is(err, gorm.ErrRecordNotFound) && result == nilpointer {
		user.Role = role
		_, err = us.Repo.CreateUser(user)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("register failed. this email have been used")
}

func (us *AuthServiceImpl) Login(user *domain.UserLogin, secret string) (*domain.PairedToken, error) {
	u, err := us.Repo.GetUserByEmail(user.Email)
	if err != nil {
		return nil, errors.New("Login failed" + err.Error())
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
		return nil, errors.New("Login failed" + err.Error())
	}
	if bcrypt.CompareHashAndPassword(hash, []byte(u.Password)) != nil {
		if err != nil {
			return nil, errors.New("Login failed" + err.Error())
		}
	}
	token, err := CreatePairedToken(user, secret)
	if user.Password != u.Password {
		if err != nil {
			return nil, errors.New("create token failed. Invalid password")
		}
	}
	return token, nil
}

func (us *AuthServiceImpl) Refresh(userID uint, token *jwt.Token, secret string) (*domain.PairedToken, error) {
	session, err := us.Ar.GetSessionByUserID(userID)
	if err != nil {
		return nil, errors.New("Refresh failed" + err.Error())
	}
	println(token.Raw)
	println(session.RefreshToken)
	if token.Raw != session.RefreshToken {
		return nil, errors.New("Refresh failed. Invalid token")
	}
	user, err := us.Repo.GetUserByID(session.UserID)
	if err != nil {
		return nil, errors.New("Refresh failed" + err.Error())
	}
	err = us.Ar.DeleteSession(userID)
	if err != nil {
		return nil, errors.New("Refresh failed" + err.Error())
	}
	newToken, err := CreatePairedToken(user, secret)
	if err != nil {
		return nil, errors.New("Refresh failed" + err.Error())
	}
	newSession := domain.Session{UserID: user.ID, RefreshToken: newToken.RefreshToken, CreatedAt: time.Now()}
	err = us.Ar.CreateSession(&newSession)
	if err != nil {
		return nil, errors.New("Refresh failed" + err.Error())
	}
	return newToken, nil
}

func (us *AuthServiceImpl) IsExpired(issuedAt time.Time) (bool, error) {
	accExp := us.Ar.Config.JwtAccessExp
	exp, err := strconv.Atoi(accExp)
	if err != nil {
		return false, errors.New("validate token failed: " + err.Error())
	}

	expiredAt := issuedAt.Add(time.Duration(exp) * time.Minute)
	if err != nil {
		return false, errors.New("Validate token failed: " + err.Error())
	}
	if time.Now().After(expiredAt) {
		return true, nil
	}
	return false, nil
}
