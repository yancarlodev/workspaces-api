package main

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/yancarlodev/workspaces-api/assets"
	"github.com/yancarlodev/workspaces-api/internal/auth"
	"github.com/yancarlodev/workspaces-api/internal/platform/app"
	"github.com/yancarlodev/workspaces-api/internal/platform/config"
	"github.com/yancarlodev/workspaces-api/internal/platform/db"
	"github.com/yancarlodev/workspaces-api/pkg/migrator"
	"github.com/yancarlodev/workspaces-api/pkg/seeder"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("available commands:\n\tmigrate\n\trun")
		return
	}

	subcommand := args[0]

	DB, err := db.Init("local.db")
	if err != nil {
		log.Fatal(err)
	}
	defer DB.Close()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		ReplaceAttr: func(groups []string, attr slog.Attr) slog.Attr {
			if attr.Key == slog.SourceKey {
				src := attr.Value.Any().(*slog.Source)
				src.File = filepath.Base(src.File)
				return slog.Any(attr.Key, src)
			}

			return attr
		},
	}))

	switch subcommand {
	case "migrate":
		fsys, err := fs.Sub(assets.Fsys, "scripts/migrations")
		if err != nil {
			logger.Error("an error has occurs when getting migrations subtree", "err", err)
		}

		m, err := migrator.New(DB, fsys, logger)
		if err != nil {
			logger.Error("an error has occurs when creating the migrator", "err", err)
		}

		if err := m.Migrate(); err != nil {
			logger.Error("an error has occurs when running the migrations", "err", err)
		}
	case "revert":
		fsys, err := fs.Sub(assets.Fsys, "scripts/migrations")
		if err != nil {
			logger.Error("an error has occurs when getting migrations subtree", "err", err)
		}

		m, err := migrator.New(DB, fsys, logger)
		if err != nil {
			logger.Error("an error has occurs when creating the migrator", "err", err)
		}

		if err := m.Revert(); err != nil {
			logger.Error("an error has occurs when reverting the migration", "err", err)
		}
	case "seed":
		fsys, err := fs.Sub(assets.Fsys, "scripts/seeds")
		if err != nil {
			logger.Error("an error has occurs when getting seeds subtree", "err", err)
		}

		s := seeder.New(DB, fsys, logger)
		if err := s.Seed(); err != nil {
			logger.Error("an error has occurs when running the seeds", "err", err)
		}
	case "run":
		startServer(DB)
	default:
		fmt.Println("available commands:\n\tmigrate\n\trun")
	}
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
