package models

import "time"

type User struct {
	ID           int       `db:"id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
}

type Subscription struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	Plan      string    `db:"plan"`
	Status    string    `db:"status"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
}

type Asset struct {
	ID     int    `db:"id"`
	Symbol string `db:"symbol"`
	Name   string `db:"name"`
	Type   string `db:"type"`
}

type PriceData struct {
	ID        int       `db:"id"`
	AssetID   int       `db:"asset_id"`
	Timestamp time.Time `db:"timestamp"`
	Open      float64   `db:"open"`
	High      float64   `db:"high"`
	Low       float64   `db:"low"`
	Close     float64   `db:"close"`
	Volume    float64   `db:"volume"`
}

type Indicator struct {
	ID        int       `db:"id"`
	AssetID   int       `db:"asset_id"`
	Timestamp time.Time `db:"timestamp"`
	RSI       float64   `db:"rsi"`
	EMAShort  float64   `db:"ema_short"`
	EMALong   float64   `db:"ema_long"`
}

type Signal struct {
	ID         int       `db:"id"`
	AssetID    int       `db:"asset_id"`
	Timestamp  time.Time `db:"timestamp"`
	SignalType string    `db:"signal_type"`
	Confidence float64   `db:"confidence"`
	Reason     string    `db:"reason"`
}
