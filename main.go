package main

import (
	"log/slog"

	"github.com/Ayikoandrew/server/api"
	"github.com/Ayikoandrew/server/database"
)

func main() {
	listenAddr := ":40000"
	store := database.NewStorage()
	if err := store.Init(); err != nil {
		slog.Error("error creating database table", "err", err)
	}
	server := api.NewServer(listenAddr, store)

	server.Run()
}
