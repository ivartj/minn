package main

import (
	"database/sql"
	_ "code.google.com/p/go-sqlite/go1/sqlite3"
)

type db interface{
	Begin() (*sql.Tx, error)
}

type dbFile struct{
	path string
	handle *sql.DB
}

type dbTemp struct{
	handle *sql.DB
}

func dbOpenFile(path string) (*dbFile) {
	return &dbFile{ path: path }
}

func (db *dbFile) Begin() (*sql.Tx, error) {

	var err error

	if db.handle == nil {
		db.handle, err = dbOpen("./" + db.path)
		if err != nil {
			return nil, err
		}
	}

	tx, err := db.handle.Begin()
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func dbOpenTemp() (*dbTemp) {
	return &dbTemp{}
}

func (db *dbTemp) Begin() (*sql.Tx, error) {

	var err error

	if db.handle == nil {
		db.handle, err = dbOpen(":memory:")
		if err != nil {
			return nil, err
		}
	}

	tx, err := db.handle.Begin()
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func dbOpen(filepath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("pragma foreign_keys = true;")
	if err != nil {
		return nil, err
	}

	return db, nil
}

