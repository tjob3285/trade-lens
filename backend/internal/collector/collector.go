package collector

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type BinancePrice struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

func StartCollector(db *sql.DB) {
	ticker := time.NewTicker(60 * time.Second)
	go func() {
		for range ticker.C {
			fetchAllPrices(db)
		}
	}()
}

func fetchAllPrices(db *sql.DB) {
	assets := []struct {
		Symbol  string
		AssetID int
	}{
		{"BTCUSDT", 1},
		{"ETHUSDT", 2},
	}

	var wg sync.WaitGroup
	for _, asset := range assets {
		wg.Add(1)
		go func(a struct {
			Symbol  string
			AssetID int
		}) {
			defer wg.Done()
			fetchAndStorePrice(db, a.Symbol, a.AssetID)
		}(asset)
	}
	wg.Wait()
}

func fetchAndStorePrice(db *sql.DB, symbol string, assetID int) {
	url := "https://api.binance.com/api/v3/ticker/price?symbol=" + symbol
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching %s: %v", symbol, err)
		return
	}
	defer resp.Body.Close()

	var data BinancePrice
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Println("Error decoding response:", err)
		return
	}

	price, err := strconv.ParseFloat(data.Price, 64)
	if err != nil {
		log.Println("Error parsing price:", err)
		return
	}

	// Insert into price_data
	_, err = db.Exec(`INSERT INTO price_data (asset_id, timestamp, open, high, low, close, volume)
                      VALUES (?, ?, ?, ?, ?, ?, ?)`,
		assetID, time.Now().Format(time.RFC3339), price, price, price, price, 0)
	if err != nil {
		log.Println("Error inserting price:", err)
		return
	}

	log.Printf("Stored BTC price: %.2f", price)
	calculateAndStoreIndicator(db, assetID)
}

func calculateAndStoreIndicator(db *sql.DB, assetID int) {
	rows, err := db.Query(`SELECT close FROM price_data WHERE asset_id = ? ORDER BY timestamp DESC LIMIT 26`, assetID)
	if err != nil {
		log.Printf("Error fetching closes for asset %d: %v", assetID, err)
		return
	}
	defer rows.Close()

	var closes []float64
	for rows.Next() {
		var close float64
		if err := rows.Scan(&close); err != nil {
			continue
		}
		closes = append(closes, close)
	}

	if len(closes) < 14 {
		log.Printf("Not enough data to calculate indicators for asset %d", assetID)
		return
	}

	for i := 0; i < len(closes)/2; i++ {
		closes[i], closes[len(closes)-1-i] = closes[len(closes)-1-i], closes[i]
	}

	emaShort := calcEMA(closes, 12)
	emaLong := calcEMA(closes, 26)

	rsi := calcRSI(closes, 14)

	_, err = db.Exec(`INSERT INTO indicators (asset_id, timestamp, rsi, ema_short, ema_long)
	VALUES (?, ?, ?, ?, ?)`,
		assetID, time.Now().Format(time.RFC3339), rsi, emaShort, emaLong)
	if err != nil {
		log.Printf("Error inserting indicators for asset %d: %v", assetID, err)
		return
	}

	log.Printf("Indicators for asset %d: RSI=%.2f, EMA12=%.2f, EMA26=%.2f", assetID, rsi, emaShort, emaLong)
}

func calcEMA(prices []float64, period int) float64 {
	if len(prices) < period {
		return prices[len(prices)-1]
	}
	k := 2.0 / float64(period+1)
	ema := prices[0]
	for i := 1; i < len(prices); i++ {
		ema = prices[i]*k + ema*(1-k)
	}
	return ema
}

func calcRSI(prices []float64, period int) float64 {
	if len(prices) < period+1 {
		return 50.0
	}

	gains := 0.0
	losses := 0.0
	for i := 1; i <= period; i++ {
		diff := prices[i] - prices[i-1]
		if diff > 0 {
			gains += diff
		} else {
			losses -= diff
		}
	}
	avgGain := gains / float64(period)
	avgLoss := losses / float64(period)
	if avgLoss == 0 {
		return 100
	}
	rs := avgGain / avgLoss
	return 100 - (100 / (1 + rs))
}
