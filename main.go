package main

import (
	"context"
	"log/slog"
	"os"

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
	if err := client.Ping(context.Background()).Err(); err != nil {
		slog.Error("Failed to connect to Redis", "error", err)
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
