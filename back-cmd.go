package main

import (
	"fmt"
	"strconv"
	"database/sql"
	"io"
	"github.com/ivartj/minn/args"
)

func init() {
	cmdRegister("back", backCmd)
}

func backCmdUsage(w io.Writer) {
	fmt.Fprintf(w, "Usage: %s back [ <card-id> ]\n", mainProgramName)

	fmt.Fprintf(w, `
Description:
  Shows the back of a card. The current card if no card ID is given.
`)
	fmt.Fprintln(w, `
Options:
  -h, --help  Prints help message.
`)
}

func backCmdArgs(cmd *cmdContext) (int, bool) {

	plainArgs := []string{}
	tok := args.NewTokenizer(cmd.Args)

	for tok.Next() {

		if tok.IsOption() {
			switch tok.Arg() {
			case "-h", "--help":
				backCmdUsage(cmd.Stdout)
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
		}
		return cardId, false
	}

	backCmdUsage(cmd.Stderr)
	cmd.Exit(1)
	return 0, false
}

func backCmd(cmd *cmdContext) {

	cardId, current := backCmdArgs(cmd)

	var err error
	if current {
		cardId, err = utilCurrentCard(cmd)
		if err != nil {
			panic(err)
		}
	}

	row := cmd.QueryRow("select back from cards where card_id = ?;", cardId)
	var back string
	err = row.Scan(&back)
	if err == sql.ErrNoRows {
		cmd.Fatalf("No card by that card ID.\n");
	}
	if err != nil {
		cmd.Fatalf("Database query error: %s.\n", err.Error())
	}

	cmd.Println(back)

	cmd.Commit()
}

