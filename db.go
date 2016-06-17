package main

import (
	"database/sql"
	_ "code.google.com/p/go-sqlite/go1/sqlite3"
)

func dbOpen(filepath string) (*sql.Tx, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("pragma foreign_keys = true;")
	if err != nil {
		return nil, err
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	return tx, nil
}

