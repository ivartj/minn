package main

import (
	"fmt"
	"io"
	"ivartj/args"
)

func init() {
	cmdRegister("add", cmdAdd, cmdAddUsage)
}

func cmdAddUsage(w io.Writer) {
	fmt.Fprintf(w, "Usage: %s add <front> <back>\n", mainProgramName)
}

func cmdAddArgs(cmd *cmdContext) (string, string) {

	plainArgs := []string{}
	tok := args.NewTokenizer(cmd.Args)

	for tok.Next() {

		if tok.IsOption() {
			switch tok.Arg() {
			case "-h", "--help":
				cmdAddUsage(cmd.Stdout)
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

	if len(plainArgs) != 2 {
		cmdAddUsage(cmd.Stderr)
		cmd.Exit(1)
	}

	return plainArgs[0], plainArgs[1]
}

func cmdAdd(cmd *cmdContext) {
	// TODO: Check against duplicate card fronts

	front, back := cmdAddArgs(cmd)

	res, err := cmd.Exec(`
		insert into
			cards (efactor, interval, front, back, entry_time)
			values (2.5, 0, ?, ?, datetime());
	`, front, back)
	if err != nil {
		cmd.Fatalf("Database error: %s.\n", err.Error())
	}

	cardId, err := res.LastInsertId()
	if err != nil {
		cmd.Fatalf("Failed to get card ID from database engine: %s.\n", err.Error())
	}

	_, err = cmd.Exec("insert into schedulings (card_id, new, schedule_time, update_efactor, update_interval) values (?, 1, datetime(), 1, 1);", cardId)
	if err != nil {
		cmd.Fatalf("Failed to schedule card: %s.\n", err.Error())
	}

	cmd.Println(cardId)

	cmd.Commit()
}

