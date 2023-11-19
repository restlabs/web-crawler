package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/we-are-discussing-rest/web-crawler/cmd/server"
	"github.com/we-are-discussing-rest/web-crawler/internal/logger"
	"github.com/we-are-discussing-rest/web-crawler/internal/repository"
	"log"
	"net/http"
	"os"
)

func main() {
	l := logger.NewLogger()
	r := repository.NewRedisRepo(context.Background(), l, &redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Username: os.Getenv("REDIS_USER"),
		Password: os.Getenv("REDIS_PW"),
	})
	s := server.NewServer(r, l)

	l.Info("listening", "PORT", os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), s))
}
