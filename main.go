package main

import (
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

	store := database.NewStorage()
	if err := store.Init(); err != nil {
		slog.Error("error creating database table", "err", err)
		os.Exit(1)
	}

	server := api.NewServer(":"+port, store)
	server.Run()
}
