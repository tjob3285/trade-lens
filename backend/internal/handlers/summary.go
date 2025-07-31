package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type SummaryResponse struct {
	Symbol     string         `json:"symbol"`
	Price      PriceResponse  `json:"price"`
	Indicators Indicator      `json:"indicators"`
	Signal     SignalResponse `json:"signal"`
}

func GetSummaryHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		symbol := r.URL.Query().Get("symbol")
		if symbol == "" {
			http.Error(w, "Missing 'symbol'", http.StatusBadRequest)
			return
		}

		// Price
		var price PriceResponse
		db.QueryRow(`
            SELECT a.symbol, p.close, p.timestamp
            FROM price_data p JOIN assets a ON p.asset_id = a.id
            WHERE a.symbol = ? ORDER BY p.timestamp DESC LIMIT 1`, symbol).Scan(&price.Symbol, &price.Price, &price.Timestamp)

		// Indicators
		var ind Indicator
		db.QueryRow(`
            SELECT i.asset_id, a.symbol, i.rsi, i.ema_short, i.ema_long
            FROM indicators i JOIN assets a ON i.asset_id = a.id
            WHERE a.symbol = ? ORDER BY i.timestamp DESC LIMIT 1`, symbol).Scan(&ind.AssetID, &ind.Symbol, &ind.RSI, &ind.EMA12, &ind.EMA26)

		// Signal
		var sig SignalResponse
		db.QueryRow(`
            SELECT a.symbol, s.signal_type, s.confidence, s.reason, s.timestamp
            FROM signals s JOIN assets a ON s.asset_id = a.id
            WHERE a.symbol = ? ORDER BY s.timestamp DESC LIMIT 1`, symbol).Scan(&sig.Symbol, &sig.SignalType, &sig.Confidence, &sig.Reason, &sig.Timestamp)

		res := SummaryResponse{Symbol: symbol, Price: price, Indicators: ind, Signal: sig}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	}
}
