package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func Connect() *sql.DB {
	db, err := sql.Open("sqlite3", "./tradelens.db")
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	createTables(db)
	db.Exec(`INSERT OR IGNORE INTO assets (id, symbol, name, type) VALUES (1, 'BTC', 'Bitcoin', 'crypto')`)
	return db
}

func createTables(db *sql.DB) {
	schema := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        email TEXT UNIQUE NOT NULL,
        password_hash TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );

    CREATE TABLE IF NOT EXISTS subscriptions (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        plan TEXT NOT NULL,
        status TEXT NOT NULL,
        expires_at DATETIME NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY(user_id) REFERENCES users(id)
    );

    CREATE TABLE IF NOT EXISTS assets (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        symbol TEXT UNIQUE NOT NULL,
        name TEXT NOT NULL,
        type TEXT NOT NULL
    );

    CREATE TABLE IF NOT EXISTS price_data (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        asset_id INTEGER NOT NULL,
        timestamp DATETIME NOT NULL,
        open REAL NOT NULL,
        high REAL NOT NULL,
        low REAL NOT NULL,
        close REAL NOT NULL,
        volume REAL NOT NULL,
        FOREIGN KEY(asset_id) REFERENCES assets(id)
    );

    CREATE TABLE IF NOT EXISTS indicators (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        asset_id INTEGER NOT NULL,
        timestamp DATETIME NOT NULL,
        rsi REAL,
        ema_short REAL,
        ema_long REAL,
        FOREIGN KEY(asset_id) REFERENCES assets(id)
    );

    CREATE TABLE IF NOT EXISTS signals (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        asset_id INTEGER NOT NULL,
        timestamp DATETIME NOT NULL,
        signal_type TEXT NOT NULL,
        confidence REAL,
        reason TEXT,
        FOREIGN KEY(asset_id) REFERENCES assets(id)
    );
    `

	_, err := db.Exec(schema)
	if err != nil {
		log.Fatal("Error creating tables:", err)
	}
	log.Println("Tables created or verified")
}
