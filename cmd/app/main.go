package main

import (
	"context"
	"github.com/vlasashk/avito-segmentation/internal/config"
	"github.com/vlasashk/avito-segmentation/internal/controller/api"
	"github.com/vlasashk/avito-segmentation/internal/model/logger"
	"github.com/vlasashk/avito-segmentation/internal/model/storage"
	"os"
)

func main() {
	cfg := config.ParseConfig()
	log := logger.InitLogger(cfg.Env)
	log.Info("Starting application")
	db, err := storage.New(context.Background())
	if err != nil {
		log.Error("Failed to initialize storage", logger.Err(err))
		os.Exit(1)
	}
	log.Info("Initialized database")
	defer db.Close()
	server := api.NewAPIServer("8080", db, log)
	//if err = db.DropDB(context.Background()); err != nil {
	//	log.Error("Failed to initialize storage", logger.Err(err))
	//	os.Exit(1)
	//}
	api.Run(log, server)
}
