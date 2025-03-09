package server

import (
	"context"
	"github.com/bobacgo/kit/examples/config"
	"log/slog"
)

const (
	KafkaServerName = "kafka"
)

type KafkaServer struct {
	conn any
}

func (k *KafkaServer) Get(name string) any {
	return k.conn
}

func (k *KafkaServer) Start(ctx context.Context) error {
	slog.Info("kafka server start", "DSN", config.Cfg().Kafka.Addr)
	return nil
}

func (k *KafkaServer) Stop(ctx context.Context) error {
	slog.Info("kafka server stop")
	return nil
}