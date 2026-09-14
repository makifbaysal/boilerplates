package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"

	rabbitmqadapter "github.com/makifbaysal/boilerplates/worker/go-kafka-worker/internal/adapter/inbound/rabbitmq"
	"github.com/makifbaysal/boilerplates/worker/go-kafka-worker/internal/adapter/outbound/memory"
	"github.com/makifbaysal/boilerplates/worker/go-kafka-worker/internal/core/job"
	"github.com/makifbaysal/boilerplates/worker/go-kafka-worker/internal/platform/config"
	"github.com/makifbaysal/boilerplates/worker/go-kafka-worker/internal/platform/logger"
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

	store := memory.NewJobStore()
	uc := job.NewService(store)

	conn, err := amqp.Dial(cfg.RabbitMQ.URL)
	if err != nil {
		log.Error("rabbitmq dial failed", "error", err)
		os.Exit(1)
	}
	defer func() { _ = conn.Close() }()

	ch, err := conn.Channel()
	if err != nil {
		log.Error("rabbitmq channel failed", "error", err)
		os.Exit(1)
	}
	defer func() { _ = ch.Close() }()

	if _, err := ch.QueueDeclare(cfg.RabbitMQ.Queue, true, false, false, false, nil); err != nil {
		log.Error("queue declare failed", "error", err)
		os.Exit(1)
	}

	deliveries, err := ch.Consume(cfg.RabbitMQ.Queue, "go-kafka-worker", false, false, false, false, nil)
	if err != nil {
		log.Error("consume failed", "error", err)
		os.Exit(1)
	}

	consumer := rabbitmqadapter.NewConsumer(uc, log)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Info("consuming", "url", cfg.RabbitMQ.URL, "queue", cfg.RabbitMQ.Queue, "env", cfg.Env)
	if err := consumer.Run(ctx, cfg.RabbitMQ.Queue, deliveries); err != nil && ctx.Err() == nil {
		log.Error("consumer stopped", "error", err)
		os.Exit(1)
	}
	log.Info("shutting down")
}
