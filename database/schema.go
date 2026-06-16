package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

const usersTable = `
CREATE TABLE IF NOT EXISTS users (
	user_id     INTEGER PRIMARY KEY AUTOINCREMENT,
    username    TEXT NOT NULL UNIQUE,
	master_hash BLOB NOT NULL,
	answer      BLOB NOT NULL,
	argon_salt  BLOB NOT NUll,
	argon_iter  INTEGER NOT NULL,
	argon_mem   INTEGER NOT NULL,
	argon_thr   INTEGER NOT NULL
);
`
const profilesTable = `
CREATE TABLE IF NOT EXISTS profiles (
	profile_id              INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id                 INTEGER NOT NULL,
	profile_hash            BLOB NOT NULL,
	encrypted_profile_name  BLOB NOT NULL,
	encrypted_pass          BLOB NOT NULL,
	created_at              DATETIME DEFAULT CURRENT_TIMESTAMP,
    update_at               DATETIME DEFAULT CURRENT_TIMESTAMP,
	unlock_at               INTEGER,
	FOREIGN KEY(user_id) REFERENCES users(user_id) ON DELETE CASCADE
);
`

func InitializeSchema(db *sql.DB) error {
	schema := usersTable + /*" " +*/ profilesTable
	_, err := db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to initialize schema: %v", err)
	}
	return nil
}

func NewDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "data.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}
	return db, nil
}

func CloseDB(db *sql.DB) error {
	err := db.Close()
	return err
}

// creates app directory in user machine and returns path to database file, error if os fails
func dbPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user directory: %v", err)
	}
	//application direcotyr
	appDir := filepath.Join(configDir, "PassMan")
	if err := os.MkdirAll(appDir, 0700); err != nil {
		return "", fmt.Errorf("failed to create data directory: %v", err)
	}
	return filepath.Join(appDir, "data.db"), nil
}
