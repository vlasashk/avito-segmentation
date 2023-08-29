package main

import (
	"github.com/vlasashk/avito-segmentation/internal/config"
	"github.com/vlasashk/avito-segmentation/internal/model"
)

func main() {
	cfg := config.ParseConfig()
	log := model.InitLogger(cfg.Env)
	log.Info("Starting application")
}
