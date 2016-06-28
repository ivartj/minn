package main

import (
	"database/sql"
	_ "code.google.com/p/go-sqlite/go1/sqlite3"
	"fmt"
	"container/list"
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

type dbMigration struct{
	from string // "" means blank state
	to string
	code string
}

// First key is to, second is from
var dbMigrations = make(map[string]map[string]*dbMigration)

func dbRegisterMigration(from, to, code string) {

	ms, ok := dbMigrations[to]
	if !ok {
		ms = make(map[string]*dbMigration)
		dbMigrations[to] = ms
	}

	ms[from] = &dbMigration{from: from, to: to, code: code}
}

func dbGetMigration(from, to string) (*dbMigration, error) {

	ms, ok := dbMigrations[to]
	if !ok {
		return nil, fmt.Errorf("No migration to %s", to)
	}

	m, ok := ms[from]
	if !ok {
		return nil, fmt.Errorf("No direct migration from %s to %s", from, to)
	}

	return m, nil
}

func dbGetMigrationPath(from, to string) ([]*dbMigration, error) {

	visited := make(map[string]string)

	success := false
	tasks := list.New()
	tasks.PushBack(to)

	for !success {

		front := tasks.Front()
		if front == nil {
			break
		}
		current := front.Value.(string)
		tasks.Remove(front)

		msToCurrent, ok := dbMigrations[current]
		if !ok {
			continue
		}

		for _, mToCurrent := range msToCurrent {

			_, alreadyVisited := visited[mToCurrent.from]
			if alreadyVisited {
				continue
			}

			visited[mToCurrent.from] = current

			if mToCurrent.from == from {
				success = true
				break
			} else {
				tasks.PushBack(mToCurrent.from)
			}

		}
	}

	if !success {
		return nil, fmt.Errorf("No path found")
	}

	path := []*dbMigration{}

	for current := from; current != to; current = visited[current] {

		m, err := dbGetMigration(current, visited[current])
		if err != nil {
			return nil, fmt.Errorf("Unexpected error retrieving migration")
		}

		path = append(path, m)
	}

	return path, nil
}

func init() {
	// TODO: Somehow use cmd.Now() instead of sqlite datetime() here

	dbRegisterMigration("", "0.1.2",`

		create table schema_changes (
			schema_version text not null,
			program_version text not null,
			change_time text not null
		);

		create table cards (
			card_id integer not null primary key,
			entry_time text not null,

			front text not null unique,
			back text not null,

			efactor float not null,
			interval integer not null,

			-- 0 new
			-- 1 relearn
			-- 2 review
			state integer not null,

			schedule_time text not null
		);

		insert into schema_changes
		values ('0.1.2', '` + mainProgramVersion + `', datetime());`)

	dbRegisterMigration("", "0.1.1",`

		create table schema_changes (
			schema_version text not null,
			program_version text not null,
			change_time text not null
		);
		
		create table cards (
			card_id integer not null primary key,
			efactor float not null,
			interval integer not null,
			front text not null unique,
			back text not null,
			entry_time text not null
		);

		create table ratings (
			time text not null,
			card_id integer not null,
			rating integer not null,
			foreign key(card_id)
				references cards(card_id)
				on delete cascade
		);

		create table schedulings (
			card_id integer not null unique,
			new integer not null,
			schedule_time text not null,
			update_efactor integer not null,
			update_interval integer not null,
			foreign key(card_id)
				references cards(card_id)
				on delete cascade
		);

		insert into schema_changes values ('0.1.1', '` + mainProgramVersion + `', datetime());`)

	dbRegisterMigration("0.1.1", "0.1.2",`

		alter table cards rename to cards_old;

		create table cards (
			card_id integer not null primary key,
			entry_time text not null,

			front text not null unique,
			back text not null,

			efactor float not null,
			interval integer not null,

			-- 0 new
			-- 1 relearn
			-- 2 review
			state integer not null,

			schedule_time text not null
		);

		insert into cards
		select
			card_id,
			entry_time,

			front,
			back,

			efactor,
			interval,

			case
				when new then 0
				when not update_interval or not update_efactor then 1
				else 2
			end as state,

			schedule_time

		from cards_old natural join schedulings;

		drop table cards_old;
		drop table ratings;
		drop table schedulings;

		insert into schema_changes values ('0.1.2', '` + mainProgramVersion + `', datetime());`)
}

