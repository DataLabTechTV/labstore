package iam

import (
	"log"
	"path/filepath"

	"github.com/IllumiKnowLabs/labstore/backend/internal/config"
	"github.com/IllumiKnowLabs/labstore/backend/internal/helper"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

type Store struct {
	Users    map[string]*User
	Groups   map[string]*Group
	Policies map[string]*Policy
}

func ensureSchema() {
	dbPath := filepath.Join(config.Storage.MetadataPath, "iam.db")

	db, err := sqlx.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalln(err)
	}
	defer helper.CloseWithErr(db, &err)

	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		access_key TEXT,
		secret_key TEXT
	);
	`

	if _, err := db.Exec(schema); err != nil {
		log.Fatalln(err)
	}
}
