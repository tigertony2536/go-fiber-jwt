package config

import (
	"os"
)

type (
	Config struct {
		DBConfig   *DBConfig
		HttpConfig *HttpConfig
	}
	DBConfig struct {
		Host     string
		Username string
		Password string
		DBName   string
		DBPort   string
		Timezone string
	}
	HttpConfig struct {
		SERVER_Port   string
		JwtSecret     string
		JwtAccessExp  string
		JwtRefreshExp string
	}
)

func NewDBConfig() *DBConfig {
	return &DBConfig{}
}
func NewHttpConfig() *HttpConfig {
	return &HttpConfig{}
}

func NewConfig() *Config {
	//DB Config
	dbConfig := NewDBConfig()
	dbConfig.Host = os.Getenv("DB_HOST")
	dbConfig.Username = os.Getenv("DB_USERNAME")
	dbConfig.Password = os.Getenv("DB_PASSWORD")
	dbConfig.DBName = os.Getenv("DB_NAME")
	dbConfig.DBPort = os.Getenv("DB_PORT")
	dbConfig.Timezone = os.Getenv("DB_TIMEZONE")

	//Http Config
	httpConifg := NewHttpConfig()
	httpConifg.SERVER_Port = os.Getenv("SERVER_PORT")
	httpConifg.JwtSecret = os.Getenv("JWT_SECRET")
	httpConifg.JwtAccessExp = os.Getenv("JWT_ACCESS_EXP")
	httpConifg.JwtAccessExp = os.Getenv("JWT_REFRESH_EXP")
	return &Config{
		DBConfig:   dbConfig,
		HttpConfig: httpConifg,
	}
}
