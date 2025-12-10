package iam

import (
	"fmt"
	"path/filepath"

	"github.com/IllumiKnowLabs/labstore/backend/internal/config"
	"github.com/IllumiKnowLabs/labstore/backend/internal/helper"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

const IAMDBFilename = "iam.db"

type Store struct {
	Users    map[string]*User
	Groups   map[string]*Group
	Policies map[string]*Policy

	readDB  *sqlx.DB
	writeDB *sqlx.DB
}

func NewStore() *Store {
	return &Store{
		Users:    make(map[string]*User),
		Groups:   make(map[string]*Group),
		Policies: make(map[string]*Policy),
	}
}

func (store *Store) open() error {
	dbPath := filepath.Join(config.Storage.MetadataPath, IAMDBFilename)

	timeoutMs := fmt.Sprint(config.IAM.DB.TimeoutMs)
	readCacheSize := fmt.Sprint(config.IAM.DB.ReadCacheSizeKiB)
	writeCacheSize := fmt.Sprint(config.IAM.DB.WriteCacheSizeKiB)

	readDSN := "file:" + dbPath +
		"?_pragma=journal_mode(WAL)" +
		"&_pragma=synchronous=NORMAL" +
		"&_pragma=temp_store=MEMORY" +
		"&_pragma=cache_size=-" + readCacheSize +
		"&_pragma=locking_mode=NORMAL" +
		"&_pragma=foreign_keys=ON" +
		"&_busy_timeout=" + timeoutMs

	writeDSN := "file:" + dbPath +
		"?_pragma=journal_mode(WAL)" +
		"&_pragma=synchronous=NORMAL" +
		"&_pragma=temp_store=MEMORY" +
		"&_pragma=cache_size=-" + writeCacheSize +
		"&_pragma=locking_mode=NORMAL" +
		"&_pragma=foreign_keys=ON" +
		"&_busy_timeout=" + timeoutMs

	readDB, err := sqlx.Open("sqlite", readDSN)
	if err != nil {
		return err
	}

	writeDB, err := sqlx.Open("sqlite", writeDSN)
	if err != nil {
		return err
	}

	readDB.SetMaxOpenConns(config.IAM.DB.MaxOpenConns)
	readDB.SetMaxIdleConns(config.IAM.DB.MaxIdleConns)

	writeDB.SetMaxOpenConns(1)
	writeDB.SetMaxIdleConns(1)

	store.readDB = readDB
	store.writeDB = writeDB

	return nil
}

func (store *Store) ensureSchema() error {
	dbPath := filepath.Join(config.Storage.MetadataPath, "iam.db")

	db, err := sqlx.Open("sqlite", dbPath)
	if err != nil {
		return err
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
		return err
	}

	return nil
}
