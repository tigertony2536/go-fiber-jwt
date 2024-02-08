package database

import (
	"fmt"
	"strconv"

	"github.com/tigertony2536/go-login/internal/config"
	"github.com/tigertony2536/go-login/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewGormDB(c *config.Config) (*gorm.DB, error) {
	port, err := strconv.Atoi(c.DBConfig.DBPort)
	if err != nil {
		return nil, err
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=%s",
		c.DBConfig.Host, c.DBConfig.Username, c.DBConfig.Password, c.DBConfig.DBName, port, c.DBConfig.Timezone)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	dbase, err := db.DB()
	if err != nil {
		return nil, err
	}
	err = dbase.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func NewUserRepositoryImpl(db *gorm.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{db: db}
}

func (ur *UserRepositoryImpl) AutoMigrate() error {
	err := ur.db.AutoMigrate(&domain.UserLogin{})
	if err != nil {
		return err
	}
	return nil
}

func (ur *UserRepositoryImpl) CreateUser(user *domain.UserLogin) error {
	var u *domain.UserLogin
	result := ur.db.Create(&u)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (ur *UserRepositoryImpl) GetUserByEmail(email string) (*domain.UserLogin, error) {
	var user *domain.UserLogin
	result := ur.db.First(&user, email)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (ur *UserRepositoryImpl) GetUsers() ([]*domain.UserLogin, error) {
	var users []*domain.UserLogin
	result := ur.db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}
func (ur *UserRepositoryImpl) UpdateUser(user *domain.UserLogin) error {
	var u domain.UserLogin
	result := ur.db.Model(&u).Updates(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (ur *UserRepositoryImpl) DeleteUser(id int) error {
	var user domain.UserLogin
	result := ur.db.Delete(&user, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
