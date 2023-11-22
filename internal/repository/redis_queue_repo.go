package repository

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/we-are-discussing-rest/web-crawler/internal/logger"
	"github.com/we-are-discussing-rest/web-crawler/internal/utils"
)

type RedisRepo struct {
	*redis.Client
	*logger.Logger
}

func NewRedisRepo(logger *logger.Logger, opts *redis.Options) *RedisRepo {
	rr := new(RedisRepo)

	rr.Logger = logger
	rr.Logger.Info("initializing redis connection")

	rdb := redis.NewClient(opts)
	rr.Client = rdb

	return rr
}

func (r *RedisRepo) CheckConnection(ctx context.Context) {
	_, err := r.Client.Ping(ctx).Result()
	if err != nil {
		r.Logger.Error("error connecting to redis", "error", err)
		panic(err)
	}

	r.Logger.Info("connected to redis")
}

func (r *RedisRepo) Insert(data string) error {
	queueName, err := utils.TrimURL(data)
	if err != nil {
		r.Logger.Error("error generating queue name", "error", err)
		fmt.Errorf("error trimming url: %v", err)
		return err
	}

	r.Client.LPush(r.Context(), queueName, data)
	r.Logger.Info("data pushed to queue", "queue", queueName, "data", data)

	return nil
}

func (r *RedisRepo) Remove(data string) error {
	return nil
}

func (r *RedisRepo) Get(data string) (string, error) {
	return "", nil
}
