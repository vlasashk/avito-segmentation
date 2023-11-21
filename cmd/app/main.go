// @title           user-segmentation api
// @version         1.0
// @description     API for dynamic user segmentation for testing new functionality

// @host      localhost:8090
// @BasePath  /

// @securityDefinitions.basic  BasicAuth

package main

import (
	"context"
	"github.com/vlasashk/user-segmentation/internal/controller/api"
	"github.com/vlasashk/user-segmentation/internal/model/logger"
	"github.com/vlasashk/user-segmentation/internal/model/storage"
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
