package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func GetPriceHistoryHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		symbol := r.URL.Query().Get("symbol")
		limit := r.URL.Query().Get("limit")
		if limit == "" {
			limit = "50"
		}

		query := `
            SELECT p.timestamp, p.close
            FROM price_data p JOIN assets a ON p.asset_id = a.id
            WHERE a.symbol = ?
            ORDER BY p.timestamp DESC LIMIT ?
        `

		rows, err := db.Query(query, symbol, limit)
		if err != nil {
			http.Error(w, "Error fetching history", 500)
			return
		}
		defer rows.Close()

		type PricePoint struct {
			Timestamp string  `json:"timestamp"`
			Price     float64 `json:"price"`
		}

		var history []PricePoint
		for rows.Next() {
			var p PricePoint
			rows.Scan(&p.Timestamp, &p.Price)
			history = append(history, p)
		}

		json.NewEncoder(w).Encode(history)
	}
}
