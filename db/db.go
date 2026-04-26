package db

import (
	"database/sql"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB(path string) error {
	var err error
	DB, err = sql.Open("sqlite", path)
	if err != nil {
		return err
	}
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE,
		password_hash TEXT,
		is_approved INTEGER DEFAULT 0,
		role TEXT DEFAULT 'user'
	);

	CREATE TABLE IF NOT EXISTS pushups (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		count INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);

	CREATE TABLE IF NOT EXISTS friendships (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		requester_id INTEGER NOT NULL,
		recipient_id INTEGER NOT NULL,
		status INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(requester_id) REFERENCES users(id),
		FOREIGN KEY(recipient_id) REFERENCES users(id),
		UNIQUE(requester_id, recipient_id)
	);`
	_, err = DB.Exec(query)
	return err
}
