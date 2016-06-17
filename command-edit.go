package main

import (
	"fmt"
	"strconv"
	"os"
	"io"
)

func commandEditUsage(w io.Writer) {
	fmt.Fprintf(w, "Usage: %s <deck> edit <card-id> <front> <back>\n", mainProgramName)
}

func commandEditArgs(cmd *cmdContext) (int, string, string) {

	plainArgs := []string{}

	for cmd.Args.Next() {

		if cmd.Args.IsOption() {
			switch cmd.Args.Arg() {
			case "-h", "--help":
				commandEditUsage(os.Stdout)
				cmd.Exit(0)
			default:
				fmt.Fprintf(os.Stderr, "Unrecognized option, '%s'.\n", cmd.Args.Arg())
				cmd.Exit(1)
			}
		} else {
			plainArgs = append(plainArgs, cmd.Args.Arg())
		}

	}

	if cmd.Args.Err() != nil {
		fmt.Fprintf(os.Stderr, "Error ccurred on parsing command-line arguments: %s.\n", cmd.Args.Err().Error())
		cmd.Exit(1)
	}

	if len(plainArgs) != 3 {
		commandEditUsage(os.Stderr)
		cmd.Exit(1)
	}

	cardId, err := strconv.Atoi(plainArgs[0])

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse card ID: %s.\n", err.Error())
		cmd.Exit(1)
	}

	return cardId, plainArgs[1], plainArgs[2]
}

func commandEdit(cmd *cmdContext) {

	cardId, front, back := commandEditArgs(cmd)

	res, err := cmd.Exec("update cards set front = ?, back = ? where card_id = ?;", front, back, cardId)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Database error: %s.\n", err.Error())
		cmd.Exit(1)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to get number of rows affected: %s.\n", err.Error())
		cmd.Exit(1)
	}

	if rowsAffected == 0 {
		fmt.Fprintf(os.Stderr, "No card by that ID.\n")
		cmd.Exit(1)
	}

	cmd.Commit()
}

