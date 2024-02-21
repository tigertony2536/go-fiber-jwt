package database

import (
	"errors"

	"github.com/tigertony2536/go-login/internal/config"
	"github.com/tigertony2536/go-login/internal/core/domain"
	"gorm.io/gorm"
)

type AuthRepositoryImpl struct {
	DB     *gorm.DB
	Config *config.HttpConfig
}

func NewAuthRepositoryImpl(db *gorm.DB, config *config.HttpConfig) *AuthRepositoryImpl {
	return &AuthRepositoryImpl{DB: db, Config: config}
}

func (ar *AuthRepositoryImpl) CreateSession(session *domain.Session) error {
	result := ar.DB.Create(session)
	if result.Error != nil {
		return errors.New("Create session failed: " + result.Error.Error())
	}
	return nil
}

func (ar *AuthRepositoryImpl) GetSessionByUserID(id uint) (*domain.Session, error) {
	var session domain.Session
	result := ar.DB.Where("user_id=?", id).First(&session)
	if result.Error != nil {
		return nil, errors.New("Get session failed: " + result.Error.Error())
	}
	return &session, nil
}

func (ar *AuthRepositoryImpl) DeleteSession(id uint) error {
	var session domain.Session
	result := ar.DB.Where("user_id=?", id).Delete(&session)
	if result.Error != nil {
		return errors.New("Delete session failed: " + result.Error.Error())
	}
	return nil
}
