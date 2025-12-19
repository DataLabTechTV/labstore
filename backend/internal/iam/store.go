package iam

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/IllumiKnowLabs/labstore/backend/internal/config"
	"github.com/IllumiKnowLabs/labstore/backend/internal/helper"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

const IAMDBFilename = "iam.db"
const defaultTTL = 15 * time.Minute

type Store struct {
	CachedUsers    map[string]*CachedUser
	CachedGroups   map[string]*CachedGroup
	CachedPolicies map[string]*CachedPolicy
	TTL            time.Duration

	readDB  *sqlx.DB
	writeCh chan<- sqlTask
}

type sqlTaskResult struct {
	sqlRes sql.Result
	err    error
}

type sqlFn func(ctx context.Context, db *sqlx.DB) sqlTaskResult

type sqlTask struct {
	uuid  string
	ctx   context.Context
	resCh chan<- sqlTaskResult
	fn    sqlFn
}

func NewStore() *Store {
	return &Store{
		CachedUsers:    make(map[string]*CachedUser),
		CachedGroups:   make(map[string]*CachedGroup),
		CachedPolicies: make(map[string]*CachedPolicy),
		TTL:            defaultTTL,
	}
}

func newSQLTask(ctx context.Context, resCh chan<- sqlTaskResult, fn sqlFn) sqlTask {
	return sqlTask{
		uuid:  uuid.NewString(),
		ctx:   ctx,
		resCh: resCh,
		fn:    fn,
	}
}

func newWriterWorker(db *sqlx.DB) chan<- sqlTask {
	slog.Debug("new writer worker", "writeChanCap", config.IAM.DB.WriteChanCap)
	taskCh := make(chan sqlTask, config.IAM.DB.WriteChanCap)

	go func() {
		for task := range taskCh {
			slog.Debug("sql write task", "uuid", task.uuid)
			sqlResp := task.fn(task.ctx, db)
			task.resCh <- sqlResp
		}
	}()

	return taskCh
}

func (store *Store) sqlExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	resCh := make(chan sqlTaskResult, 1)

	store.writeCh <- newSQLTask(ctx, resCh,
		func(ctx context.Context, db *sqlx.DB) sqlTaskResult {
			res, err := db.ExecContext(ctx, query, args...)
			return sqlTaskResult{sqlRes: res, err: err}
		},
	)

	res := <-resCh

	return res.sqlRes, res.err
}

func (store *Store) sqlNamedExecContext(ctx context.Context, query string, arg any) (sql.Result, error) {
	resCh := make(chan sqlTaskResult, 1)

	store.writeCh <- newSQLTask(ctx, resCh,
		func(ctx context.Context, db *sqlx.DB) sqlTaskResult {
			res, err := db.NamedExecContext(ctx, query, arg)
			return sqlTaskResult{sqlRes: res, err: err}
		},
	)

	res := <-resCh

	return res.sqlRes, res.err
}

func (store *Store) open() error {
	dbPath := filepath.Join(config.Storage.MetadataPath, IAMDBFilename)

	f, err := os.OpenFile(dbPath, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	helper.CloseWithErr(f, &err)

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
	store.writeCh = newWriterWorker(writeDB)

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
		user_id TEXT PRIMARY KEY,
		name TEXT UNIQUE,
		arn TEXT UNIQUE,
		access_key TEXT,
		secret_key BLOB
	);

	CREATE TABLE IF NOT EXISTS groups (
		group_id TEXT PRIMARY KEY,
		name TEXT UNIQUE,
		arn TEXT UNIQUE
	);

	CREATE TABLE IF NOT EXISTS policies (
		policy_id TEXT PRIMARY KEY,
		name TEXT UNIQUE,
		arn TEXT UNIQUE,

		document JSON NOT NULL,

		created_at DATETIME NOT NULL DEFAULT (CURRENT_TIMESTAMP),
		updated_at DATETIME NOT NULL DEFAULT (CURRENT_TIMESTAMP)
	);

	CREATE TABLE IF NOT EXISTS group_users (
		group_id TEXT,
		user_id TEXT,
		PRIMARY KEY (group_id, user_id),
		FOREIGN KEY(group_id)
			REFERENCES groups(group_id)
			ON DELETE CASCADE
			ON UPDATE CASCADE,
		FOREIGN KEY(user_id)
			REFERENCES users(user_id)
			ON DELETE CASCADE
			ON UPDATE CASCADE
	);

	CREATE TRIGGER IF NOT EXISTS policies_update_trigger
	AFTER UPDATE ON policies
	FOR EACH ROW
	BEGIN
		UPDATE policies
		SET updated_at = CURRENT_TIMESTAMP
		WHERE policy_id = OLD.policy_id;
	END;

	CREATE TABLE IF NOT EXISTS user_policies (
		user_id TEXT,
		policy_id TEXT,
		PRIMARY KEY (user_id, policy_id),
		FOREIGN KEY(user_id)
			REFERENCES users(user_id)
			ON DELETE CASCADE
			ON UPDATE CASCADE,
		FOREIGN KEY(policy_id)
			REFERENCES policies(policy_id)
			ON DELETE CASCADE
			ON UPDATE CASCADE
	);

	CREATE TABLE IF NOT EXISTS group_policies (
		group_id TEXT,
		policy_id TEXT,
		PRIMARY KEY (group_id, policy_id),
		FOREIGN KEY(group_id)
			REFERENCES groups(group_id)
			ON DELETE CASCADE
			ON UPDATE CASCADE,
		FOREIGN KEY(policy_id)
			REFERENCES policies(policy_id)
			ON DELETE CASCADE
			ON UPDATE CASCADE
	);
	`

	if _, err := db.Exec(schema); err != nil {
		return err
	}

	return nil
}
