package main

import (
	"fmt"
	"strconv"
	"io"
	"github.com/ivartj/minn/args"
)

func init() {
	cmdRegister("remove", removeCmd)
}

func removeCmdUsage(w io.Writer) {
	fmt.Fprintf(w, "Usage: %s remove [ <card-id> ]\n", mainProgramName)

	fmt.Fprintf(w, `
Description:
  Removes the card by the given card ID.

Options:
  -h, --help  Prints help message.

`)

}

func removeCmdArgs(cmd *cmdContext) (int, bool) {

	plainArgs := []string{}
	tok := args.NewTokenizer(cmd.Args)

	for tok.Next() {

		if tok.IsOption() {
			switch tok.Arg() {
			case "-h", "--help":
				removeCmdUsage(cmd.Stdout)
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

	removeCmdUsage(cmd.Stderr)
	cmd.Exit(1)
	return 0, false
}

func removeCmd(cmd *cmdContext) {

	cardId, current := removeCmdArgs(cmd)

	var err error
	if current {
		cardId, err = utilCurrentCard(cmd)
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

