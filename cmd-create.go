package main

import (
	"fmt"
	"io"
	"ivartj/args"
)

func init() {
	cmdRegister("create", cmdCreate, cmdCreateUsage)
}

func cmdCreateUsage(w io.Writer) {
	fmt.Fprintf(w, "Usage: %s create\n", mainProgramName)
}

func cmdCreateArgs(cmd *cmdContext) {

	tok := args.NewTokenizer(cmd.Args)

	for tok.Next() {

		if tok.IsOption() {

			switch tok.Arg() {
			case "-h", "--help":
				cmdCreateUsage(cmd.Stdout)
				cmd.Exit(0)
			default:
				cmd.Fatalf("Unrecognized option, '%s'.\n", tok.Arg())
			}
				
		} else {
			cmdCreateUsage(cmd.Stderr)
			cmd.Exit(1)
		}
	}

	if tok.Err() != nil {
		cmd.Fatalf("Error occurred on processing command-line arguments: %s.\n", tok.Err().Error())
	}
}

func cmdCreate(cmd *cmdContext) {

	cmdCreateArgs(cmd)

	_, err := cmd.Exec(`

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

		insert into schema_changes values (?, ?, ?);

	`, mainSchemaVersion, mainProgramVersion, cmd.Now().Format(utilTimeFormat))

	if err != nil {
		cmd.Fatalf("Error on creating database: %s.\n", err.Error())
	}

	cmd.Commit()
}

