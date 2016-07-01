package main

import (
	"ivartj/args"
	"io"
	"fmt"
)

func init() {
	cmdRegister("migrate", migrateCmd)
}

func migrateCmdUsage(w io.Writer) {
	fmt.Fprintf(w, "Usage: %s migrate [ -t <to-schema-version> ]\n", mainProgramName)
}

func migrateCmdArgs(cmd *cmdContext) (schemaVersion string, fromBlank bool) {

	tok := args.NewTokenizer(cmd.Args)
	schemaVersion = mainSchemaVersion

	for tok.Next() {

		if tok.IsOption() {

			switch tok.Arg() {

			case "-h", "--help":
				migrateCmdUsage(cmd.Stdout)
				cmd.Exit(0)

			case "-t", "--to":
				param, err := tok.TakeParameter()
				if err != nil {
					cmd.Fatalf("Failed to read parameter to '%s' option: %s.\n", err.Error())
				}
				schemaVersion = param

			case "-b", "--from-blank":
				fromBlank = true

			default:
				cmd.Fatalf("Unrecognized option, '%s'.\n", tok.Arg())

			}

		} else {
			migrateCmdUsage(cmd.Stderr)
			cmd.Exit(1)
		}

	}

	if tok.Err() != nil {
		cmd.Fatalf("Error occurred on parsing command-line arguments: %s.\n", tok.Err().Error())
	}

	return schemaVersion, fromBlank
}

func migrateCmd(cmd *cmdContext) {
	schemaVersion, fromBlank := migrateCmdArgs(cmd)

	var currentSchemaVersion string
	if fromBlank {
		currentSchemaVersion = ""
	} else {
		row := cmd.QueryRow(`select schema_version from schema_changes order by change_time desc limit 1;`)
		err := row.Scan(&currentSchemaVersion)
		if err != nil {
			cmd.Fatalf("Error occurred on getting the current schema version: %s.\n", err.Error())
		}
	}

	path, err := dbGetMigrationPath(currentSchemaVersion, schemaVersion)
	if err != nil {
		cmd.Fatalf("Failed to get a migration path: %s.\n", err.Error())
	}

	for _, m := range path {

		fmt.Fprintf(cmd.Stderr, "%s -> %s\n", m.from, m.to)
		_, err := cmd.Exec(m.code)
		if err != nil {
			cmd.Fatalf("Failed to apply the migration %s -> %s: %s.\n", m.from, m.to, err.Error())
		}
	}

	cmd.Commit()
}

