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
