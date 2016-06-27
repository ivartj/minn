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

		insert into schema_changes values (?, ?, datetime());

	`, mainSchemaVersion, mainProgramVersion)

	if err != nil {
		cmd.Fatalf("Error on creating database: %s.\n", err.Error())
	}

	cmd.Commit()
}

