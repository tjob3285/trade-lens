package main

import (
	"trade-lens/internal/collector"
	"trade-lens/internal/config"
	"trade-lens/internal/db"
	"trade-lens/internal/server"
)

func main() {
	cfg := config.Load()
	database := db.Connect()
	defer database.Close()

	collector.StartCollector(database)

	server.Start(cfg.Port, database)
}
