package database

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/tigertony2536/go-login/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewGormDB(c *config.Config) (*gorm.DB, error) {
	port, err := strconv.Atoi(c.DBConfig.DBPort)
	if err != nil {
		return nil, err
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=%s",
		c.DBConfig.Host, c.DBConfig.Username, c.DBConfig.Password, c.DBConfig.DBName, port, c.DBConfig.Timezone)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: NewLogger(),
	})
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

func NewLogger() logger.Interface {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  true,        // Disable color
		},
	)
	return newLogger
}
