package main

import (
	"log/slog"
	"net/http"
	_ "net/http/pprof"
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

	go func() {
		slog.Info("Starting pprof server on :6060")
		if err := http.ListenAndServe(":6060", nil); err != nil {
			slog.Error("pprof server failed", "err", err)
		}
	}()

	database.InitRedis()

	store := database.NewStorage()
	if err := store.Init(); err != nil {
		slog.Error("error creating database table", "err", err)
		os.Exit(1)
	}

	server := api.NewServer(":"+port, store)
	server.Run()
	server.StartTokenCleanup(24 * time.Hour)
}
