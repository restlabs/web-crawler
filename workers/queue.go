package main

import (
	"context"
	"github.com/we-are-discussing-rest/web-crawler/workers/repository"
	"log/slog"
)

func MapQueues(queue *repository.RedisRepo, logger *slog.Logger, ctx context.Context) []string {
	keys := queue.GetAllKeys(ctx)
	logger.Info("got mapping for all queues")
	return keys
}
