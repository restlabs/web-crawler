package main

import (
	poolparty "github.com/we-are-discussing-rest/pool-party"
	"github.com/we-are-discussing-rest/web-crawler/workers/repository"
	"github.com/we-are-discussing-rest/web-crawler/workers/utils"
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	store := repository.NewSqliteRepo(utils.Lookup("SQLITE_URL", "data.db"), logger)

	p := poolparty.NewPool(5)
	c := NewCrawler(CrawlerOpts{
		urls:       []string{"https://example.com"},
		store:      store,
		workerPool: p,
		logger:     logger,
	})

	ct := NewCrawler(CrawlerOpts{
		urls:       []string{"https://en.wikipedia.org/wiki/B-tree"},
		store:      store,
		workerPool: p,
		logger:     logger,
	})

	p.Start()
	c.Crawl()
	ct.Crawl()

	p.Stop()
}
