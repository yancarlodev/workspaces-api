package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/yancarlodev/workspaces-api/internal/platform/config"
)

func main() {
	server := gin.Default()
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	if err := server.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal(err)
	}
}
