package main

import (
	"fmt"
	"strconv"
	"io"
	"ivartj/args"
)

func init() {
	cmdRegister("remove", cmdRemove, cmdRemoveUsage)
}

func cmdRemoveUsage(w io.Writer) {
	fmt.Fprintf(w, "Usage: %s remove [ <card-id> ]\n", mainProgramName)
}

func cmdRemoveArgs(cmd *cmdContext) (int, bool) {

	plainArgs := []string{}
	tok := args.NewTokenizer(cmd.Args)

	for tok.Next() {

		if tok.IsOption() {
			switch tok.Arg() {
			case "-h", "--help":
				cmdRemoveUsage(cmd.Stdout)
				cmd.Exit(0)
			default:
				cmd.Fatalf("Unrecognized option: %s.\n", tok.Arg())
			}
		} else {
			plainArgs = append(plainArgs, tok.Arg())
		}
	}

	if tok.Err() != nil {
		cmd.Fatalf("Error on processing command-line arguments: %s.\n", tok.Err().Error())
	}

	switch(len(plainArgs)) {
	case 0:
		return 0, true
	case 1:
		cardId, err := strconv.Atoi(plainArgs[0])
		if err != nil {
			cmd.Fatalf("Failed to parse Card ID: %s.\n", err.Error())
			cmd.Exit(1)
		}
		return cardId, false
	}

	cmdRemoveUsage(cmd.Stderr)
	cmd.Exit(1)
	return 0, false
}

func cmdRemove(cmd *cmdContext) {

	cardId, current := cmdRemoveArgs(cmd)

	var err error
	if current {
		cardId, err = sm2CurrentCard(cmd.DB())
		if err != nil {
			panic(err)
		}
	}

	res, err := cmd.Exec("delete from cards where card_id = ?;", cardId)
	if err != nil {
		cmd.Fatalf("Database error: %s.\n", err.Error())
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		cmd.Fatalf("Unable to get number of rows affected: %s.\n", err.Error())
	}

	if rowsAffected == 0 {
		cmd.Fatalf("No card by that card ID.\n")
	}

	cmd.Commit()
}

