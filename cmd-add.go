package main

import (
	"fmt"
	"os"
	"io"
	"ivartj/args"
)

func init() {
	cmdRegister("add", cmdAdd, cmdAddUsage)
}

func cmdAddUsage(w io.Writer) {
	fmt.Fprintf(w, "Usage: %s <deck> add <front> <back>\n", mainProgramName)
}

func cmdAddArgs(cmd *cmdContext) (string, string) {

	plainArgs := []string{}
	tok := args.NewTokenizer(cmd.Args)

	for tok.Next() {

		if tok.IsOption() {
			switch tok.Arg() {
			case "-h", "--help":
				cmdAddUsage(os.Stdout)
				cmd.Exit(0)
			default:
				fmt.Fprintf(os.Stderr, "Unrecognized option, '%s'.\n", tok.Arg())
				cmd.Exit(1)
			}
		} else {
			plainArgs = append(plainArgs, tok.Arg())
		}

	}

	if tok.Err() != nil {
		fmt.Fprintf(os.Stderr, "Error ccurred on parsing command-line arguments: %s.\n", tok.Err().Error())
		cmd.Exit(1)
	}

	if len(plainArgs) != 2 {
		cmdAddUsage(os.Stderr)
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
		fmt.Fprintf(os.Stderr, "Database error: %s.\n", err.Error())
		cmd.Exit(1)
	}

	cardId, err := res.LastInsertId()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get card ID from database engine: %s.\n", err.Error())
		cmd.Exit(1)
	}

	_, err = cmd.Exec("insert into schedulings (card_id, new, schedule_time, update_efactor, update_interval) values (?, 1, datetime(), 1, 1);", cardId)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to schedule card: %s.\n", err.Error())
		cmd.Exit(1)
	}

	fmt.Printf("%d\n", cardId)

	cmd.Commit()
}

