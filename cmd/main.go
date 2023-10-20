package main

import (
	"auth/api"
	"auth/api/handler"
	"auth/config"
	"auth/pkg/logger"
	"auth/storage/postgres"
	"context"
	"fmt"
)

func main() {
	cfg := config.Load()
	log := logger.NewLogger("blog-project", logger.LevelInfo)
	strg, err := postgres.NewStorage(context.Background(), cfg)
	if err != nil {
		return
	}

	h := handler.NewHandler(strg, log)

	r := api.NewServer(h)
	r.Run(fmt.Sprintf(":%s", cfg.Port))
}
