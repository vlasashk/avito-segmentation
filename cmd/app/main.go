package main

import (
	"fmt"
	"github.com/vlasashk/avito-segmentation/internal/config"
)

func main() {
	cfg := config.ParseConfig()
	fmt.Println(cfg)
}
