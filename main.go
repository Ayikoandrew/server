package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/Ayikoandrew/server/api"
	"github.com/Ayikoandrew/server/database"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	store := database.NewStorage()
	if err := store.Init(); err != nil {
		slog.Error("error creating database table", "err", err)
	}
	server := api.NewServer(":"+port, store)
	log.Printf("Server starting on port %s...", port)

	server.Run()
}
