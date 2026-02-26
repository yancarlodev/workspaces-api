package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/yancarlodev/workspaces-api/internal/auth"
	"github.com/yancarlodev/workspaces-api/internal/platform/app"
	"github.com/yancarlodev/workspaces-api/internal/platform/config"
)

func main() {
	server := gin.Default()
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	application := app.New()

	registerRoutes(server, application)

	if err := server.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal(err)
	}
}

func registerRoutes(server *gin.Engine, application *app.App) {
	auth.RegisterRoutes(server, application.AuthHandler)
}
