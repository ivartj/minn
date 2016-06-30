package main

import (
	"fmt"
	"os"
	"database/sql"
	"io"
	"time"
)

type cmdContext struct{
	tx *sql.Tx
	db db
	Args []string
	Stdout io.Writer
	Stderr io.Writer
	Stdin io.Reader
	offset time.Duration // For manipulation in tests only
	MaxRelearnBacklog int
}

type cmdExit int

var cmdList = []cmdListItem{}

type cmdListItem struct{
	name string
	fn func(*cmdContext)
}

func cmdRegister(name string, fn func(*cmdContext)) {
	cmdList = append(cmdList, cmdListItem{name, fn})
}

func cmdNewContext(db db) *cmdContext {
	return &cmdContext{
		db: db,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin: os.Stdin,
		offset: 0,
		MaxRelearnBacklog: mainConfMaxRelearnBacklog,
	}
}

func (cmd *cmdContext) Printf(format string, args... interface{}) {
	_, err := fmt.Fprintf(cmd.Stdout, format, args...)
	if err != nil {
		cmd.Fatalf("cmd.Printf error: %s.\n", err.Error())
	}
}

func (cmd *cmdContext) Print(args... interface{}) {
	_, err := fmt.Fprint(cmd.Stdout, args...)
	if err != nil {
		cmd.Fatalf("cmd.Print error: %s.\n", err.Error())
	}
}

func (cmd *cmdContext) Println(args... interface{}) {
	_, err := fmt.Fprintln(cmd.Stdout, args...)
	if err != nil {
		cmd.Fatalf("cmd.Print error: %s.\n", err.Error())
	}
}

func (cmd *cmdContext) Fatalf(format string, args... interface{}) {
	fmt.Fprintf(cmd.Stderr, format, args...)
	cmd.Exit(1)
}

func (cmd *cmdContext) Fatal(args... interface{}) {
	fmt.Fprint(cmd.Stderr, args...)
	cmd.Exit(1)
}

func (cmd *cmdContext) Fatalln(args... interface{}) {
	fmt.Fprintln(cmd.Stderr, args...)
	cmd.Exit(1)
}

func (cmd *cmdContext) Warnf(format string, args... interface{}) {
	fmt.Fprintf(cmd.Stderr, format, args...)
}

func (cmd *cmdContext) Warn(args... interface{}) {
	fmt.Fprint(cmd.Stderr, args...)
}

func (cmd *cmdContext) Warnln(args... interface{}) {
	fmt.Fprintln(cmd.Stderr, args...)
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

func (cmd *cmdContext) Run(iargv... interface{}) (status int) {

	argv := make([]string, len(iargv))
	for i, v := range iargv {
		argv[i] = fmt.Sprint(v)
	}
	cmd.Args = argv

	if len(argv) < 1 {
		fmt.Fprintf(cmd.Stderr, "No command given.\n")
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
		fmt.Fprintf(cmd.Stderr, "Unrecognized command, '%s'.\n", argv[0])
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
			fmt.Fprintf(cmd.Stdout, "Panic occurred: %s.\n", err.Error())
			status = 1
		}
		cmd.Rollback()
	}()

	cmd.tx = nil
	fn(cmd)

	return status
}

func (cmd *cmdContext) ForwardTime(d time.Duration) {
	cmd.offset += d
}

func (cmd *cmdContext) Now() time.Time {
	return time.Now().Add(cmd.offset)
}

