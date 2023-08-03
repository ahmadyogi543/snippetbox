package main

import (
	"database/sql"
	"time"
)

func openDB(driverName string, dsn string) (*sql.DB, error) {
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func formatHumanReadableDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}
