package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/we-are-discussing-rest/web-crawler/workers/utils"
	"log/slog"
)

var (
	NoMessageError = errors.New("no messages in the queue")
)

type RedisRepo struct {
	*redis.Client
	logger *slog.Logger
}

func NewRedisRepo(logger *slog.Logger, opts *redis.Options) *RedisRepo {
	rr := new(RedisRepo)

	rr.logger = logger
	rr.logger.Info("initializing redis connection")

	rdb := redis.NewClient(opts)
	rr.Client = rdb

	return rr
}

func (r *RedisRepo) CheckConnection(ctx context.Context) {
	_, err := r.Client.Ping(ctx).Result()
	if err != nil {
		r.logger.Error("error connecting to redis", "error", err)
		panic(err)
	}

	r.logger.Info("connected to redis")
}

func (r *RedisRepo) Insert(data string) error {
	queueName, err := utils.TrimURL(data)
	if err != nil {
		r.logger.Error("error generating queue name", "error", err)
		return fmt.Errorf("error trimming url: %v", err)
	}

	r.Client.LPush(r.Context(), queueName, data)
	r.logger.Info("data pushed to queue", "queue", queueName, "data", data)

	return nil
}

func (r *RedisRepo) Remove(data string) error {
	return nil
}

func (r *RedisRepo) Get(data string) (string, error) {
	r.logger.Info("dequeue", "key", data)
	val := r.Client.RPop(r.Context(), data)
	if val.Val() == "" {
		return "", NoMessageError
	}

	return val.Val(), nil
}

func (r *RedisRepo) GetAllMessages(data string, ctx context.Context) ([]string, error) {
	r.logger.Info("getting queue len", "key", data)
	ql := r.Client.LLen(ctx, data)
	val := r.Client.RPopCount(ctx, data, int(ql.Val()))
	if len(val.Val()) == 0 {
		return nil, NoMessageError
	}

	return val.Val(), nil
}

func (r *RedisRepo) GetAllKeys(ctx context.Context) []string {
	ret := r.Client.Keys(ctx, "*")

	return ret.Val()
}
