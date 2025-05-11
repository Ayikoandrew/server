package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/Ayikoandrew/server/api"
	"github.com/Ayikoandrew/server/database"
)

func main() {
	slog.Info("Starting server...")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	database.InitRedis()

	client := database.GetRedisClient()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		slog.Error("Failed to connect to Redis", "error", err, "url", os.Getenv("REDIS_URL"))
		os.Exit(1)
	}
	slog.Info("Successfully connected to Redis")

	store := database.NewStorage()
	if err := store.Init(); err != nil {
		slog.Error("error creating database table", "err", err)
	}
	server := api.NewServer(":"+port, store)

	server.Run()
}
