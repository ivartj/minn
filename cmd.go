package main

import (
	"fmt"
	"os"
	"database/sql"
)

type cmdContext struct{
	tx *sql.Tx
	deckfilepath string
	Args []string
}

type cmdExit int

func cmdNewContext(deckfilepath string, argv []string) *cmdContext {
	return &cmdContext{
		deckfilepath: deckfilepath,
		Args: argv,
	}
}

func (cmd *cmdContext) initDB() {
	if cmd.tx == nil {
		var err error
		cmd.tx, err = dbOpen(cmd.deckfilepath)
		if err != nil {
			panic(fmt.Errorf("Database error: %s", err.Error()))
		}
	}
}

func (cmd *cmdContext) DB() *sql.Tx {
	cmd.initDB()
	return cmd.tx
}

func (cmd *cmdContext) Query(q string, params... interface{}) (*sql.Rows, error) {
	cmd.initDB()

	return cmd.tx.Query(q, params...)
}

func (cmd *cmdContext) QueryRow(q string, params... interface{}) (*sql.Row) {
	cmd.initDB()

	return cmd.tx.QueryRow(q, params...)
}

func (cmd *cmdContext) Exec(q string, params... interface{}) (sql.Result, error) {
	cmd.initDB()

	return cmd.tx.Exec(q, params...)
}

func (cmd *cmdContext) Exit(status int) {
	panic(cmdExit(status))
}

func (cmd *cmdContext) Rollback() {
	if cmd.tx == nil {
		return
	}

	err := cmd.tx.Rollback()
	if err == sql.ErrTxDone {
		return
	}
	if err != nil {
		panic(err)
	}
}

func (cmd *cmdContext) Commit() {
	if cmd.tx == nil {
		return
	}

	err := cmd.tx.Commit()
	if err == sql.ErrTxDone {
		return
	}
	if err != nil {
		panic(err)
	}
}

func (cmd *cmdContext) Execute(fn func(*cmdContext)) (status int) {
	status = 0
	defer func() {
		x := recover()
		switch x.(type) {
		case cmdExit:
			status = int(x.(cmdExit))
		case error:
			err := x.(error)
			fmt.Fprintf(os.Stderr, "Panic occurred: %s.\n", err.Error())
			status = 1
		}
		cmd.Rollback()
	}()
	fn(cmd)
	return status
}

