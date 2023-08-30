package main

import (
	"context"
	"fmt"
	"github.com/vlasashk/avito-segmentation/internal/config"
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
	defer db.Close()
	fmt.Println(db)
	//var greeting string
	//err = dbpool.QueryRow(ctx, "select 'Hello, world!'").Scan(&greeting)
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	//	os.Exit(1)
	//}
	//
	//fmt.Println(greeting)
}
