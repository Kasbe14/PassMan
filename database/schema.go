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
	user_id        INTEGER PRIMARY KEY AUTOINCREMENT,
    username       TEXT NOT NULL UNIQUE,
	master_hash    TEXT NOT NULL,
    encrypt_salt   BLOB NOT NULL,
	WrappedKeyPass BLOB NOT NULL, 
	WrappedKeyRec  BLOB NOT NULL 
);
`
const profilesTable = `
CREATE TABLE IF NOT EXISTS profiles (
	profile_id              INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id                 INTEGER NOT NULL,
	pro_hash                TEXT NOT NULL,
	enc_pro_name            BLOB NOT NULL,
	enc_pass                BLOB NOT NULL,
	created_at              INTEGER NOT NULL,
    update_at               INTEGER NOT NULL,
    lck                     BOOL,
	unlock_at               INTEGER,
	FOREIGN KEY(user_id) REFERENCES users(user_id) ON DELETE CASCADE
);
`
//TODO ;adding prama forreing key on to the databse and 

//   //TODO at startup logic for database register a hook to run on EVERY new connection
    // sqlite.RegisterConnectionHook(func(conn sqlite.ExecQuerierContext, dsn string) error {
    //     _, err := conn.ExecContext(context.Background(), "PRAGMA foreign_keys = ON", nil)
    //     return err
    // })

func InitializeSchema(db *sql.DB) error {
	schema := usersTable + /*" " +*/ profilesTable
	_, err := db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to initialize schema: %v", err)
	}
	return nil
}

//creates new databse  at userConfigdiretory
func NewDatabase() (*sql.DB, error) {
     // userDbPath, err := dbPath(); 
     // if err != nil {
     //    return nil, fmt.Errorf("failed to open database: %v", err)
    // }
    db, err := sql.Open("sqlite", /*userDbPath*/"./data.db")
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
