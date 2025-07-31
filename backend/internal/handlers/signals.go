package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Signal struct {
	ID         int       `json:"id"`
	AssetID    int       `json:"asset_id"`
	Timestamp  time.Time `json:"timestamp"`
	SignalType string    `json:"signal_type"`
	Confidence float64   `json:"confidence"`
	Reason     string    `json:"reason"`
}

func GetSignalsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(`SELECT id, asset_id, timestamp, signal_type, confidence, reason FROM signals`)
		if err != nil {
			log.Println("Error fetching signals:", err)
			http.Error(w, "Fetch failed", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var signals []Signal
		for rows.Next() {
			var s Signal
			var ts string
			err := rows.Scan(&s.ID, &s.AssetID, &ts, &s.SignalType, &s.Confidence, &s.Reason)
			if err != nil {
				log.Println("Error scanning row:", err)
				continue
			}
			s.Timestamp, _ = time.Parse(time.RFC3339Nano, ts)
			signals = append(signals, s)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(signals)
	}
}

type SignalResponse struct {
	Symbol     string  `json:"symbol"`
	SignalType string  `json:"signal_type"`
	Confidence float64 `json:"confidence"`
	Reason     string  `json:"reason"`
	Timestamp  string  `json:"timestamp"`
}

func GetLatestSignalHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		symbol := r.URL.Query().Get("symbol")
		if symbol == "" {
			http.Error(w, "Missing 'symbol'", http.StatusBadRequest)
			return
		}

		query := `
            SELECT a.symbol, s.signal_type, s.confidence, s.reason, s.timestamp
            FROM signals s
            JOIN assets a ON s.asset_id = a.id
            WHERE a.symbol = ?
            ORDER BY s.timestamp DESC LIMIT 1
        `

		var res SignalResponse
		err := db.QueryRow(query, symbol).Scan(&res.Symbol, &res.SignalType, &res.Confidence, &res.Reason, &res.Timestamp)
		if err != nil {
			log.Println("Error fetching signal:", err)
			http.Error(w, "Signal not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	}
}
