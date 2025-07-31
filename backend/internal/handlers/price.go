package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type PriceResponse struct {
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	Timestamp string  `json:"timestamp"`
}

func GetPriceHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		symbol := r.URL.Query().Get("symbol")
		if symbol == "" {
			http.Error(w, "Missing 'symbol'", http.StatusBadRequest)
			return
		}

		query := `
            SELECT a.symbol, p.close, p.timestamp
            FROM price_data p
            JOIN assets a ON p.asset_id = a.id
            WHERE a.symbol = ?
            ORDER BY p.timestamp DESC LIMIT 1
        `

		var res PriceResponse
		err := db.QueryRow(query, symbol).Scan(&res.Symbol, &res.Price, &res.Timestamp)
		if err != nil {
			log.Println("Error fetching price:", err)
			http.Error(w, "Price not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	}
}
