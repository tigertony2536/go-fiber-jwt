package database

import (
	"errors"

	"github.com/tigertony2536/go-login/internal/core/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepositoryImpl(db *gorm.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{db: db}
}

func (ur *UserRepositoryImpl) CreateUser(user *domain.UserLogin) (*gorm.DB, error) {
	result := ur.db.Table("user_logins").Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return result, nil
}

func (ur *UserRepositoryImpl) GetUserByID(id uint) (*domain.UserLogin, error) {
	var user domain.UserLogin
	result := ur.db.Table("user_logins").First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (ur *UserRepositoryImpl) GetUserByEmail(email string) (*domain.UserLogin, error) {
	var user domain.UserLogin
	result := ur.db.Table("user_logins").Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (ur *UserRepositoryImpl) GetUsers() (*[]domain.UserLogin, error) {
	var users []domain.UserLogin
	result := ur.db.Table("user_logins").Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return &users, nil
}

func (ur *UserRepositoryImpl) UpdateUser(email string, user *domain.UserLogin) error {
	u, err := ur.GetUserByEmail(email)
	if err != nil {
		return errors.New("Not found user with this ID" + err.Error())
	}
	result := ur.db.Table("user_logins").Model(&u).Updates(map[string]interface{}{"Email": user.Email, "Password": user.Password})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (ur *UserRepositoryImpl) DeleteUser(email string) error {
	var user domain.UserLogin
	result := ur.db.Clauses(clause.Returning{}).Where("email = ?", email).Delete(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
