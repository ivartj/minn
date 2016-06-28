package main

import (
	"fmt"
	"strconv"
	"database/sql"
	"io"
	"ivartj/args"
)

func init() {
	cmdRegister("back", cmdBack, cmdBackUsage)
}

func cmdBackUsage(w io.Writer) {
	fmt.Fprintf(w, "Usage: %s back [ <card-id> ]\n", mainProgramName)
}

func cmdBackArgs(cmd *cmdContext) (int, bool) {

	plainArgs := []string{}
	tok := args.NewTokenizer(cmd.Args)

	for tok.Next() {

		if tok.IsOption() {
			switch tok.Arg() {
			case "-h", "--help":
				cmdBackUsage(cmd.Stdout)
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

	cmdBackUsage(cmd.Stderr)
	cmd.Exit(1)
	return 0, false
}

func cmdBack(cmd *cmdContext) {

	cardId, current := cmdBackArgs(cmd)

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

