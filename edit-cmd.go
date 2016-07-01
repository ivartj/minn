package main

import (
	"fmt"
	"strconv"
	"io"
	"ivartj/args"
)

func init() {
	cmdRegister("edit", editCmd)
}

func editCmdUsage(w io.Writer) {
	fmt.Fprintf(w, "Usage: %s edit <card-id> <front> <back>\n", mainProgramName)

	fmt.Fprintf(w, `
Description:
  Edits front and back of the card by the given card ID.
`)

	fmt.Fprintln(w, `
Options:
  -h, --help  Prints help message.
`)
}

func editCmdArgs(cmd *cmdContext) (int, string, string) {

	plainArgs := []string{}
	tok := args.NewTokenizer(cmd.Args)

	for tok.Next() {

		if tok.IsOption() {
			switch tok.Arg() {
			case "-h", "--help":
				editCmdUsage(cmd.Stdout)
				cmd.Exit(0)
			default:
				cmd.Fatalf("Unrecognized option, '%s'.\n", tok.Arg())
			}
		} else {
			plainArgs = append(plainArgs, tok.Arg())
		}

	}

	if tok.Err() != nil {
		cmd.Fatalf("Error ccurred on parsing command-line arguments: %s.\n", tok.Err().Error())
	}

	if len(plainArgs) != 3 {
		editCmdUsage(cmd.Stderr)
		cmd.Exit(1)
	}

	cardId, err := strconv.Atoi(plainArgs[0])

	if err != nil {
		cmd.Fatalf("Failed to parse card ID: %s.\n", err.Error())
	}

	return cardId, plainArgs[1], plainArgs[2]
}

func editCmd(cmd *cmdContext) {

	cardId, front, back := editCmdArgs(cmd)

	res, err := cmd.Exec("update cards set front = ?, back = ? where card_id = ?;", front, back, cardId)
	if err != nil {
		cmd.Fatalf("Database error: %s.\n", err.Error())
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		cmd.Fatalf("Unable to get number of rows affected: %s.\n", err.Error())
	}

	if rowsAffected == 0 {
		cmd.Fatalf("No card by that ID.\n")
	}

	cmd.Commit()
}

