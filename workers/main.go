package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	poolparty "github.com/we-are-discussing-rest/pool-party"
	"github.com/we-are-discussing-rest/web-crawler/workers/repository"
	"github.com/we-are-discussing-rest/web-crawler/workers/utils"
	"log/slog"
	"os"
	"sync/atomic"
)

func main() {
	ctx := context.Background()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	store := repository.NewSqliteRepo(utils.Lookup("SQLITE_URL", "data.db"), logger)
	queue := repository.NewRedisRepo(logger, &redis.Options{
		Addr:     utils.Lookup("REDIS_HOST", "localhost:6379"),
		Username: utils.Lookup("REDIS_USER", ""),
		Password: utils.Lookup("REDIS_PW", ""),
	})
	p := poolparty.NewPool(20)
	p.Start()

	m := MapQueues(queue, logger, ctx)

	for _, v := range m {
		am, err := queue.GetAllMessages(v, ctx)
		if err != nil {
			return
		}

		NewCrawler(CrawlerOpts{
			currDepth:  &atomic.Uint64{},
			urls:       am,
			store:      store,
			queue:      queue,
			workerPool: p,
			logger:     logger,
		}).Crawl()
	}

	defer p.Stop()
}
