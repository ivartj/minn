package main

import (
	"fmt"
	"strconv"
	"database/sql"
	"ivartj/args"
	"io"
)

func init() {
	cmdRegister("front", frontCmd)
}

func frontCmdUsage(w io.Writer) {
	fmt.Fprintf(w, "Usage: %s front [ <card-id> ]\n", mainProgramName)
}

func frontCmdArgs(cmd *cmdContext) (int, bool) {

	plainArgs := []string{}
	tok := args.NewTokenizer(cmd.Args)

	for tok.Next() {

		if tok.IsOption() {
			switch tok.Arg() {
			case "-h", "--help":
				frontCmdUsage(cmd.Stdout)
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
	default:
		frontCmdUsage(cmd.Stderr)
		cmd.Exit(1)
	}

	// unreachable
	return 0, false
}

func frontCmd(cmd *cmdContext) {

	cardId, current := frontCmdArgs(cmd)

	var err error
	if current {
		cardId, err = utilCurrentCard(cmd)
		if err != nil {
			panic(err)
		}
	}

	row := cmd.QueryRow("select front from cards where card_id = ?;", cardId)
	var front string
	err = row.Scan(&front)
	if err == sql.ErrNoRows {
		cmd.Fatalf("No card by that card ID.\n");
	}
	if err != nil {
		cmd.Fatalf("Database query error: %s.\n", err.Error())
	}

	cmd.Println(front)

	cmd.Commit()
}

