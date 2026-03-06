package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/yancarlodev/workspaces-api/internal/auth"
	"github.com/yancarlodev/workspaces-api/internal/platform/app"
	"github.com/yancarlodev/workspaces-api/internal/platform/config"
	"github.com/yancarlodev/workspaces-api/internal/platform/db"
)

func main() {
	DB, err := db.Init("local.db")
	if err != nil {
		log.Fatal(err)
	}
	defer DB.Close()

	startServer(DB)
}

func startServer(DB *sql.DB) {
	server := gin.Default()
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	application := app.New(DB)

	registerRoutes(server, application)

	if err := server.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal(err)
	}
}

func registerRoutes(server *gin.Engine, application *app.App) {
	auth.RegisterRoutes(server, application.AuthHandler)
}
