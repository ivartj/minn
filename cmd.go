package main

import (
	"fmt"
	"database/sql"
	"github.com/ivartj/norske-irc-kanaler.com/args"
)

type cmdContext struct{
	tx *sql.Tx
	deckfilepath string
	Args *args.Tokenizer
}

type cmdExit int

func cmdNewContext(deckfilepath string, tok *args.Tokenizer) *cmdContext {
	return &cmdContext{
		deckfilepath: deckfilepath,
		Args: tok,
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


