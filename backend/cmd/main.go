package main

import (
	model "backend/cmd/api/models"
	"backend/cmd/config"
	"backend/cmd/database"
	server "backend/cmd/server"
	"log"
	"github.com/joho/godotenv"
)

func main () {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	cfg := config.NewConfig()
	if err = cfg.Validate(); err != nil {
		log.Fatal(err)
	}

	db := database.NewDatabaseConnection(cfg.DBConf.DBAddr)
	urlStore := model.NewUrlStore(db)

	server.Start(cfg, urlStore)
}

