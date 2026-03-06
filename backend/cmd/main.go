package main

import (
	model "backend/cmd/api/models"
	"backend/cmd/config"
	"backend/cmd/database"
	server "backend/cmd/server"
	"log/slog"

	"github.com/joho/godotenv"
)

func main () {
	if err := godotenv.Load(); err != nil {
    	slog.Info("[WARNING] Missing .env file, using environment variables")
	}

	cfg := config.NewConfig()
	if err := cfg.Validate(); err != nil {
		slog.Info("[WARNING]", "error", err)
	}

	db := database.NewDatabaseConnection(cfg.DBConf.DBAddr)
	urlStore := model.NewUrlStore(db)

	server.Start(cfg, urlStore)
}

