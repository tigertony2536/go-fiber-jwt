package database

import (
	"database/sql"

	"github.com/tigertony2536/go-login/internal/core/domain"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Use this database For Unit test
func NewTestDB() (*gorm.DB, *sql.DB, error) {
	gorm, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}
	db, err := gorm.DB()
	if err != nil {
		return nil, nil, err
	}
	err = gorm.AutoMigrate(&domain.UserLogin{}, &domain.Session{})
	if err != nil {
		return nil, nil, err
	}
	return gorm, db, nil
}
