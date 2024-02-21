package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/tigertony2536/go-login/internal/adaptor/database"
	"github.com/tigertony2536/go-login/internal/adaptor/router"
	"github.com/tigertony2536/go-login/internal/config"
	"github.com/tigertony2536/go-login/internal/core"
	"github.com/tigertony2536/go-login/internal/core/domain"
)

func init() {
	envPath := "..\\.env"
	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("Error loading .env file from %s ", envPath)
	}
}

func main() {
	config := config.NewConfig()
	// db, err := database.NewGormDB(config)
	gorm, err := database.NewGormDB(config)
	if err != nil {
		log.Fatal("can not init gorm: ", err)
	}
	db, err := gorm.DB()
	if err != nil {
		log.Fatal("can not call gorm's db: ", err)
	}
	defer db.Close()
	err = gorm.AutoMigrate(&domain.UserLogin{}, &domain.Session{})
	if err != nil {
		log.Fatal("migration fail: ", err)
	}
	userRepo := database.NewUserRepositoryImpl(gorm)
	authRepository := database.NewAuthRepositoryImpl(gorm, config.HttpConfig)
	authService := core.NewAuthService(userRepo, authRepository)
	authHandler := router.NewAuthHandler(authService, authRepository)
	userHandler := router.NewUserHandler(userRepo)
	rt := router.NewRouter(authHandler, userHandler, config)
	err = rt.Serve(config.HttpConfig)
	if err != nil {
		log.Fatal(err)
	}

}
