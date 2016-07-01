package main

import (
	"fmt"
	"io"
	"ivartj/args"
)

func init() {
	cmdRegister("create", createCmd)
}

func createCmdUsage(w io.Writer) {
	fmt.Fprintf(w, "Usage: %s create\n", mainProgramName)
}

func createCmdArgs(cmd *cmdContext) {

	tok := args.NewTokenizer(cmd.Args)

	for tok.Next() {

		if tok.IsOption() {

			switch tok.Arg() {
			case "-h", "--help":
				createCmdUsage(cmd.Stdout)
				cmd.Exit(0)
			default:
				cmd.Fatalf("Unrecognized option, '%s'.\n", tok.Arg())
			}
				
		} else {
			createCmdUsage(cmd.Stderr)
			cmd.Exit(1)
		}
	}

	if tok.Err() != nil {
		cmd.Fatalf("Error occurred on processing command-line arguments: %s.\n", tok.Err().Error())
	}
}

func createCmd(cmd *cmdContext) {

	createCmdArgs(cmd)

	path, err := dbGetMigrationPath("", mainSchemaVersion)
	if err != nil {
		cmd.Fatalf("Did not find migration path to current schema version: %s.\n", err.Error())
	}

	for _, m := range path {

		_, err = cmd.Exec(m.code)
		if err != nil {
			cmd.Fatalf("Failed to apply migration %s -> %s.\n", m.from, m.to)
		}

	}

	cmd.Commit()
}

