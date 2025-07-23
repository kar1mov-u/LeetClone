package main

import (
	"context"
	"log"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/kar1mov-u/LeetClone/internal/api"
	database "github.com/kar1mov-u/LeetClone/internal/repo"
	"github.com/kar1mov-u/LeetClone/internal/services"
)

func main() {
	_ = godotenv.Load()
	dbConfig, err := database.NewConfig()
	if err != nil {
		log.Fatalf("error loading database config: %v", err)
	}

	err = database.Migrate(context.Background(), *dbConfig, "file:///home/karimov/Projects/LeetClone/backend/migrations")
	if err != nil {
		log.Fatalf("error on migrations: %v", err)
	}

	pool, err := database.NewConn(context.TODO(), *dbConfig)
	if err != nil {
		log.Fatal(err)
	}
	userRepo := database.NewUserRepo(pool)
	userService := services.NewUserService(userRepo)

	api := api.New(userService)
	api.Start()

}
