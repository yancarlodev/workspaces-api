package app

import "github.com/yancarlodev/workspaces-api/internal/auth"

type App struct {
	AuthHandler *auth.Handler
}

func New() *App {
	authHandler := auth.NewHandler()

	return &App{
		AuthHandler: authHandler,
	}
}
