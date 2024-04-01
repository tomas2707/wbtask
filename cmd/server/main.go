package main

import (
	"github.com/tomas2707/wbtask"
	"github.com/tomas2707/wbtask/repository/gorm"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg := gorm.DBConfig{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USER"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		Password: os.Getenv("DB_PASSWORD"),
	}

	db, err := gorm.InitDB(cfg)
	if err != nil {
		panic(err)
	}

	userRepo := gorm.NewUserRepository(db)
	userService := wbtask.NewService(logger, userRepo)

	http.HandleFunc("POST /save", userService.SaveUserHandler)
	http.HandleFunc("GET /{id}", userService.GetUserHandler)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
