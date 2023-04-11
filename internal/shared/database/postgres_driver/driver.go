package postgres_driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}
var dbInitialized bool

const maxOpenDBConnections = 10
const maxIdleDBConn = 5
const maxDBLifetime = 5 * time.Minute

func ConnectSQL(dsn string) (*DB, error) {
	if dbInitialized {
		return dbConn, nil
	}
	db, err := NewDatabase(dsn)

	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(maxOpenDBConnections)
	db.SetConnMaxIdleTime(maxIdleDBConn)
	db.SetConnMaxLifetime(maxDBLifetime)

	dbConn.SQL = db
	dbInitialized = true

	return dbConn, nil
}

func NewDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
