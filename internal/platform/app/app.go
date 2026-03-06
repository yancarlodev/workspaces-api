package app

import (
	"database/sql"

	"github.com/yancarlodev/workspaces-api/internal/auth"
)

type App struct {
	AuthHandler *auth.Handler
}

func New(DB *sql.DB) *App {
	authHandler := auth.NewHandler()

	return &App{
		AuthHandler: authHandler,
	}
}
