package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/we-are-discussing-rest/web-crawler/cmd/server"
	"github.com/we-are-discussing-rest/web-crawler/logger"
	"github.com/we-are-discussing-rest/web-crawler/repository"
	"github.com/we-are-discussing-rest/web-crawler/utils"
	"log"
	"net/http"
	"os"
)

func main() {
	l := logger.NewLogger()

	r := repository.NewRedisRepo(l, &redis.Options{
		Addr:     utils.Lookup("REDIS_HOST", "localhost:6379"),
		Username: utils.Lookup("REDIS_USER", ""),
		Password: utils.Lookup("REDIS_PW", ""),
	})
	r.CheckConnection(context.Background())

	s := server.NewServer(r, l)

	l.Info("listening", "PORT", os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), s))
}
