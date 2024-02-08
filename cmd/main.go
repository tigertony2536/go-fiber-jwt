package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/tigertony2536/go-login/internal/adaptor/database"
	"github.com/tigertony2536/go-login/internal/adaptor/router"
	"github.com/tigertony2536/go-login/internal/config"
)

func main() {
	err := godotenv.Load("..\\.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	config := config.NewConfig()
	db, err := database.NewGormDB(config)
	if err != nil {
		log.Fatal(err)
	}
	userRepo := database.NewUserRepositoryImpl(db)
	authHandler := router.NewAuthHandler(userRepo)
	rt := router.NewRouter(authHandler, config)
	err = rt.Serve(config.HttpConfig)
	if err != nil {
		log.Fatal(err)
	}
}
