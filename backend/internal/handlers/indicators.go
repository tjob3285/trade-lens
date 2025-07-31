package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type Indicator struct {
	AssetID int     `json:"asset_id"`
	Symbol  string  `json:"symbol"`
	RSI     float64 `json:"rsi"`
	EMA12   float64 `json:"ema_short"`
	EMA26   float64 `json:"ema_long"`
}

func GetIndicatorsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		symbol := r.URL.Query().Get("symbol")
		if symbol == "" {
			http.Error(w, "Missing symbol query param", http.StatusBadRequest)
		}

		query := `
			SELECT i.asset_id, a.symbol, i.rsi, i.ema_short, i.ema_long
			FROM indicators i
			JOIN assets a ON i.asset_id = a.id
			WHERE a.symbol = ?
			ORDER BY i.timestamp DESC LIMIT 1
		`

		var ind Indicator
		err := db.QueryRow(query, symbol).Scan(&ind.AssetID, &ind.Symbol, &ind.RSI, &ind.EMA12, &ind.EMA26)
		if err != nil {
			log.Printf("Error fetching indicators for %s: %v", symbol, err)
			http.Error(w, "Indicators not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ind)
	}
}
