package core

import "database/sql"

//provides database operations to the tui
type VaultService struct {
	db *sql.DB
}
