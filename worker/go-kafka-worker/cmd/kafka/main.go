package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"

	kafkaadapter "github.com/makifbaysal/boilerplates/worker/go-kafka-worker/internal/adapter/inbound/kafka"
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

	client, err := kgo.NewClient(
		kgo.SeedBrokers(cfg.Kafka.Brokers...),
		kgo.ConsumeTopics(cfg.Kafka.Topic),
		kgo.ConsumerGroup(cfg.Kafka.GroupID),
		kgo.DisableAutoCommit(),
	)
	if err != nil {
		log.Error("kafka client init failed", "error", err)
		os.Exit(1)
	}
	defer client.Close()

	consumer := kafkaadapter.NewConsumer(client, uc, log).
		WithRateLimit(cfg.RateLimit.PerSecond, cfg.RateLimit.Burst,
			time.Duration(cfg.RateLimit.WaitSeconds)*time.Second)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Info("consuming", "brokers", cfg.Kafka.Brokers, "topic", cfg.Kafka.Topic, "group_id", cfg.Kafka.GroupID, "env", cfg.Env)
	if err := consumer.Run(ctx); err != nil && ctx.Err() == nil {
		log.Error("consumer stopped", "error", err)
		os.Exit(1)
	}
	log.Info("shutting down")
}
