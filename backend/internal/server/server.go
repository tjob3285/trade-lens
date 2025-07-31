package server

import (
	"database/sql"
	"log"
	"net/http"
	"trade-lens/internal/handlers"
)

func Start(port string, db *sql.DB) {
	http.HandleFunc("/health", handlers.HealthHandler)
	http.HandleFunc("/signals", handlers.GetSignalsHandler(db))
	http.HandleFunc("/indicators", handlers.GetIndicatorsHandler(db))
	http.HandleFunc("/dashboard", handlers.DashboardHandler(db))
	http.HandleFunc("/price", handlers.GetPriceHandler(db))
	http.HandleFunc("/signal", handlers.GetLatestSignalHandler(db))
	http.HandleFunc("/summary", handlers.GetSummaryHandler(db))
	http.HandleFunc("/prices/history", handlers.GetPriceHistoryHandler(db))

	http.HandleFunc("/auth/register", handlers.RegisterHandler(db))
	http.HandleFunc("/auth/login", handlers.LoginHandler(db))

	log.Println("Starting server on port", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
