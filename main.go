package main

import (
	"fmt"
	"github.com/we-are-discussing-rest/web-crawler/cmd/server"
	"github.com/we-are-discussing-rest/web-crawler/internal/logger"
	"github.com/we-are-discussing-rest/web-crawler/internal/repository"
	"log"
	"net/http"
	"os"
)

func main() {
	l := logger.NewLogger()
	r := repository.NewRedisRepo()
	s := server.NewServer(r, l)

	l.Info("listening", "PORT", os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), s))
}
