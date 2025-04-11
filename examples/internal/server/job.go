package server

import (
	"context"
	"log/slog"
)

const (
	JobServerName = "job"
)

type JobServer struct {
	task []string
}

func (k *JobServer) Get() any {
	return k
}

func (k *JobServer) Start(ctx context.Context) error {
	slog.Info("job server start")
	return nil
}

func (k *JobServer) Stop(ctx context.Context) error {
	slog.Info("job server stop")
	return nil
}

func (k *JobServer) AddTask(ctx context.Context, task string) error {
	k.AddTask(ctx, task)
	return nil
}
