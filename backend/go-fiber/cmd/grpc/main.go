package main

import (
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	grpcadapter "github.com/makifbaysal/boilerplates/backend/go-fiber/internal/adapter/inbound/grpc"
	"github.com/makifbaysal/boilerplates/backend/go-fiber/internal/adapter/inbound/grpc/taskpb"
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

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(logger.UnaryServerInterceptor(log)))
	taskpb.RegisterTaskServiceServer(grpcServer, grpcadapter.NewServer(uc))

	healthSrv := health.NewServer()
	healthSrv.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(grpcServer, healthSrv)

	reflection.Register(grpcServer)

	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "9090"
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Error("listen failed", "error", err)
		os.Exit(1)
	}

	log.Info("listening", "port", port, "env", cfg.Env, "transport", "grpc")
	if err := grpcServer.Serve(lis); err != nil {
		log.Error("server stopped", "error", err)
		os.Exit(1)
	}
}
