package main

import (
	"os"

	gqlhandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"

	graphqladapter "github.com/makifbaysal/boilerplates/backend/go-fiber/internal/adapter/inbound/graphql"
	"github.com/makifbaysal/boilerplates/backend/go-fiber/internal/adapter/inbound/graphql/generated"
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
	if err := cfgLoader.Watch(log, func(newCfg *config.Config) {
		level.Set(logger.ParseLevel(newCfg.Log.Level))
	}); err != nil {
		log.Error("config watch setup failed, continuing without live reload", "error", err)
	}

	repo := memory.NewTaskRepository()
	uc := task.NewService(repo)

	schema := generated.NewExecutableSchema(generated.Config{Resolvers: graphqladapter.NewResolver(uc)})
	gqlServer := gqlhandler.NewDefaultServer(schema)

	app := fiber.New(fiber.Config{AppName: "go-fiber-boilerplate-graphql"})
	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(logger.FiberMiddleware(log))

	app.Get("/healthz", func(c fiber.Ctx) error { return c.JSON(fiber.Map{"status": "ok"}) })
	app.All("/query", adaptor.HTTPHandler(gqlServer))
	app.Get("/", adaptor.HTTPHandlerFunc(playground.Handler("GraphQL Playground", "/query")))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	log.Info("listening", "port", port, "env", cfg.Env, "transport", "graphql")
	if err := app.Listen(":" + port); err != nil {
		log.Error("server stopped", "error", err)
		os.Exit(1)
	}
}
