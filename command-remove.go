package main

import (
	"fmt"
	"strconv"
	"os"
	"io"
)

func commandRemoveUsage(w io.Writer) {
	fmt.Fprintf(w, "Usage: %s <deck> remove [ <card-id> ]\n", mainProgramName)
}

func commandRemoveArgs(cmd *cmdContext) (int, bool) {

	plainArgs := []string{}

	for cmd.Args.Next() {

		if cmd.Args.IsOption() {
			switch cmd.Args.Arg() {
			case "-h", "--help":
				commandRemoveUsage(os.Stdout)
				cmd.Exit(0)
			default:
				fmt.Fprintf(os.Stderr, "Unrecognized option: %s.\n", cmd.Args.Arg())
				cmd.Exit(1)
			}
		} else {
			plainArgs = append(plainArgs, cmd.Args.Arg())
		}
	}

	if cmd.Args.Err() != nil {
		fmt.Fprintf(os.Stderr, "Error on processing command-line arguments: %s.\n", cmd.Args.Err().Error())
		cmd.Exit(1)
	}

	switch(len(plainArgs)) {
	case 0:
		return 0, true
	case 1:
		cardId, err := strconv.Atoi(plainArgs[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse Card ID: %s.\n", err.Error())
			cmd.Exit(1)
		}
		return cardId, false
	}

	commandRemoveUsage(os.Stderr)
	cmd.Exit(1)
	return 0, false
}

func commandRemove(cmd *cmdContext) {

	cardId, current := commandRemoveArgs(cmd)

	var err error
	if current {
		cardId, err = sm2CurrentCard(cmd.DB())
		if err != nil {
			panic(err)
		}
	}

	res, err := cmd.Exec("delete from cards where card_id = ?;", cardId)
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
		fmt.Fprintf(os.Stderr, "No card by that card ID.\n")
		cmd.Exit(1)
	}

	cmd.Commit()
}

