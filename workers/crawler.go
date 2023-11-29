package main

import (
	"errors"
	"fmt"
	poolparty "github.com/we-are-discussing-rest/pool-party"
	"github.com/we-are-discussing-rest/web-crawler/workers/internal"
	"github.com/we-are-discussing-rest/web-crawler/workers/repository"
	"log/slog"
)

type Crawler struct {
	opts CrawlerOpts
}

type CrawlerOpts struct {
	urls       []string
	store      repository.Repository
	workerPool *poolparty.Pool
	logger     *slog.Logger
}

func NewCrawler(opts CrawlerOpts) *Crawler {
	return &Crawler{
		opts: CrawlerOpts{
			urls:       opts.urls,
			store:      opts.store,
			workerPool: opts.workerPool,
			logger:     opts.logger,
		},
	}
}

func (c *Crawler) Crawl() {
	for _, url := range c.opts.urls {
		c.opts.workerPool.Send(func() {
			c.opts.logger.Info("starting a crawl", "url", url)

			err := internal.ResolveDns(url)
			if errors.Is(err, internal.ErrorIpCannotBeResolved) {
				c.opts.logger.Warn("DNS could not resolve for url", "url", url)
				return
			}

			rawHtml, err := internal.DownloadRawHtml(url)
			if err != nil {
				c.opts.logger.Error("error downloading raw HTML", "error", err)
				return
			}

			validateErr := internal.ValidateHtmlContent(rawHtml)
			if validateErr != nil {
				c.opts.logger.Error("error validating html content", "error", validateErr, "url", url)
				return
			}

			pageErr := internal.CheckContent(c.opts.store, rawHtml)
			if errors.Is(pageErr, internal.ContentErrorDuplicateHash) {
				c.opts.logger.Info("duplicate hash value, page already scraped", "url", url)
				return
			} else if pageErr != nil {
				c.opts.logger.Error("error checking content", "error", pageErr)
				return
			}

			links, linkErr := internal.ExtractHtmlLinks(rawHtml)
			if linkErr != nil {
				c.opts.logger.Error("error extracting links", "error", linkErr, "url", url)
				return
			}

			for _, link := range links {
				fmt.Println(link)
			}
		})
	}
}
