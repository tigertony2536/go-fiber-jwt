package config_test

import (
	"log"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/tigertony2536/go-login/internal/config"
)

func init() {
	err := godotenv.Load("D:\\dev\\go\\src\\03-side-projects\\go-login\\.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func TestNewConfig(t *testing.T) {
	const (
		JWT_SECRET      = "abc123abc"
		JWT_ACCESS_EXP  = "15"
		JWT_REFRESH_EXP = "10080"
		DB_HOST         = "localhost"
		DB_USERNAME     = "myuser"
		DB_PASSWORD     = "mypassword"
		DB_NAME         = "mydatabase"
		DB_PORT         = "5432"
		DB_TIMEZONE     = "Asia/Shanghai"
		DB_file         = "test.db"
		SERVER_PORT     = ":8000"
	)
	cfg := config.NewConfig()
	assert.Equal(t, JWT_SECRET, cfg.HttpConfig.JwtSecret)
}
