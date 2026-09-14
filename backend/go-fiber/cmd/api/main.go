package main

import (
	"os"
	"time"

	"github.com/makifbaysal/boilerplates/backend/go-fiber/internal/adapter/inbound/rest"
	"github.com/makifbaysal/boilerplates/backend/go-fiber/internal/adapter/outbound/memory"
	"github.com/makifbaysal/boilerplates/backend/go-fiber/internal/core/task"
	"github.com/makifbaysal/boilerplates/backend/go-fiber/internal/platform/config"
	"github.com/makifbaysal/boilerplates/backend/go-fiber/internal/platform/logger"
)

func main() {
	configDir := os.Getenv("CONFIG_DIR")
	if configDir == "" {
		configDir = "configs"
	}

	cfgLoader, err := config.New(configDir)
	if err != nil {
		panic(err)
	}
	cfg := cfgLoader.Current()

	log, level := logger.NewDynamic(logger.Config{Level: cfg.Log.Level, Format: cfg.Log.Format})

	// Dynamic config: editing configs/config.<env>.yaml while the
	// process is running updates the log level in place, no restart.
	// See .ai/common-tasks.md for adding more fields to this flow.
	if err := cfgLoader.Watch(log, func(newCfg *config.Config) {
		level.Set(logger.ParseLevel(newCfg.Log.Level))
	}); err != nil {
		log.Error("config watch setup failed, continuing without live reload", "error", err)
	}

	repo := memory.NewTaskRepository()
	uc := task.NewService(repo)
	app := rest.NewWithRateLimit(uc, rest.RateLimit{
		Disabled: cfg.RateLimit.Disabled,
		Max:      cfg.RateLimit.Max,
		Window:   time.Duration(cfg.RateLimit.WindowSeconds) * time.Second,
	}, logger.FiberMiddleware(log))

	log.Info("listening", "port", cfg.Server.Port, "env", cfg.Env)
	if err := app.Listen(":" + cfg.Server.Port); err != nil {
		log.Error("server stopped", "error", err)
		os.Exit(1)
	}
}
