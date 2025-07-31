package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
)

type DashboardData struct {
	Assets []AssetData `json:"assets"`
}

type AssetData struct {
	Symbol     string  `json:"symbol"`
	Price      float64 `json:"price"`
	RSI        float64 `json:"rsi"`
	EMA12      float64 `json:"ema_short"`
	EMA26      float64 `json:"ema_long"`
	Signal     string  `json:"signal"`
	Confidence float64 `json:"confidence"`
}

func DashboardHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		symbols := []string{"BTC", "ETH"}
		var assets []AssetData

		for _, symbol := range symbols {
			asset := AssetData{Symbol: symbol}

			err := db.QueryRow(`
			SELECT p.close
			FROM price_data p
			JOIN assets a ON p.asset_id = a.id
			WHERE a.symbol = ?
			ORDER BY p.timestamp DESC LIMIT 1`, symbol+"USDT").Scan(&asset.Price)
			if err != nil {
				log.Println("Error getting price:", err)
				continue
			}

			err = db.QueryRow(`
			SELECT rsi, ema_short, ema_long
			FROM indicators i
			JOIN assets a ON i.asset_id = a.id
			WHERE a.symbol = ?
			ORDER BY i.timestamp DESC LIMIT 1`, symbol+"USDT").Scan(&asset.RSI, &asset.EMA12, &asset.EMA26)
			if err != nil {
				log.Println("Error getting indicators:", err)
			}

			err = db.QueryRow(`
			SELECT signal_type, confidence
			FROM signals s
			JOIN assets a ON s.asset_id = a.id
			WHERE a.symbol = ?
			ORDER BY s.timestamp DESC LIMIT 1`, symbol+"USDT").Scan(&asset.Signal, &asset.Confidence)
			if err != nil {
				log.Println("Error getting signal:", err)
			}

			assets = append(assets, asset)
		}

		tmpl := `
        <html>
        <head>
            <title>TradeLens Dashboard</title>
            <style>
                body { font-family: Arial; margin: 20px; }
                table { width: 50%; border-collapse: collapse; }
                th, td { border: 1px solid #ddd; padding: 8px; text-align: center; }
                th { background: #f4f4f4; }
            </style>
			<script>
				setTimeout(function() {
					window.location.reload();
				}, 10000); // refresh every 10 seconds
			</script>
        </head>
        <body>
            <h1>TradeLens Dashboard</h1>
            <table>
                <tr>
                    <th>Asset</th>
                    <th>Price</th>
                    <th>RSI</th>
                    <th>EMA12</th>
                    <th>EMA26</th>
                    <th>Signal</th>
                    <th>Confidence</th>
                </tr>
                {{range .Assets}}
                <tr>
                    <td>{{.Symbol}}</td>
                    <td>{{printf "%.2f" .Price}}</td>
                    <td>{{printf "%.2f" .RSI}}</td>
                    <td>{{printf "%.2f" .EMA12}}</td>
                    <td>{{printf "%.2f" .EMA26}}</td>
                    <td>{{.Signal}}</td>
                    <td>{{printf "%.0f%%" .Confidence}}</td>
                </tr>
                {{end}}
            </table>
        </body>
        </html>
        `

		t, err := template.New("dashboard").Parse(tmpl)
		if err != nil {
			http.Error(w, "Template error", http.StatusInternalServerError)
			return
		}

		t.Execute(w, DashboardData{Assets: assets})
	}
}
