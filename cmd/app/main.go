package main

import (
	"context"
	"github.com/vlasashk/avito-segmentation/internal/controller/api"
	"github.com/vlasashk/avito-segmentation/internal/model/logger"
	"github.com/vlasashk/avito-segmentation/internal/model/storage"
	"os"
)

func main() {
	log := logger.InitLogger()
	log.Info("Starting application")
	db, err := storage.New(context.Background())
	if err != nil {
		log.Error("Failed to initialize storage", logger.Err(err))
		os.Exit(1)
	}
	log.Info("Database successfully initialized")
	defer db.Close()
	server := api.NewAPIServer(os.Getenv("PORT"), db, log)
	api.Run(log, server)
}
