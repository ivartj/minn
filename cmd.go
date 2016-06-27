package main

import (
	"fmt"
	"os"
	"database/sql"
	"io"
)

type cmdContext struct{
	tx *sql.Tx
	db db
	Args []string
}

type cmdExit int

var cmdList = []cmdListItem{}

type cmdListItem struct{
	name string
	fn func(*cmdContext)
	usage func(io.Writer)
}

func cmdRegister(name string, fn func(*cmdContext), usage func(io.Writer)) {
	cmdList = append(cmdList, cmdListItem{name, fn, usage})
}

func cmdNewContext(db db) *cmdContext {
	return &cmdContext{
		db: db,
	}
}

func (cmd *cmdContext) initDB() {
	if cmd.tx == nil {
		var err error
		cmd.tx, err = cmd.db.Begin()
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

func (cmd *cmdContext) Run(argv []string) (status int) {

	if len(argv) < 1 {
		fmt.Fprintf(os.Stderr, "No command given.\n")
		return 1
	}

	var fn func(*cmdContext) = nil
	for _, v := range cmdList {
		if v.name == argv[0] {
			fn = v.fn
			break
		}
	}

	if fn == nil {
		fmt.Fprintf(os.Stderr, "Unrecognized command, '%s'.\n", argv[0])
		return 1
	}

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

	cmd.tx = nil
	cmd.Args = argv
	fn(cmd)

	return status
}

