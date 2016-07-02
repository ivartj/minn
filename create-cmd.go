package main

import (
	"fmt"
	"io"
	"github.com/ivartj/minn/args"
)

func init() {
	cmdRegister("create", createCmd)
}

func createCmdUsage(w io.Writer) {
	fmt.Fprintf(w, "Usage: %s create\n", mainProgramName)

	fmt.Fprintf(w, `
Description:
  Initializes deck database. As with other subcommands, the path of the
  database is by default at:

   %s

  This can be changed with the --deck option before this subcommand.
`, mainConfDeckPath)

	fmt.Fprintln(w, `
Options:
  -h, --help  Prints help message.
`)

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

