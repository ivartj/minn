package main

import (
	"fmt"
	"os"
	"io"
)

func commandAddUsage(w io.Writer) {
	fmt.Fprintf(w, "Usage: %s <deck> add <front> <back>\n", mainProgramName)
}

func commandAddArgs(cmd *cmdContext) (string, string) {

	plainArgs := []string{}

	for cmd.Args.Next() {

		if cmd.Args.IsOption() {
			switch cmd.Args.Arg() {
			case "-h", "--help":
				commandAddUsage(os.Stdout)
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

	if len(plainArgs) != 2 {
		commandAddUsage(os.Stderr)
		cmd.Exit(1)
	}

	return plainArgs[0], plainArgs[1]
}

func commandAdd(cmd *cmdContext) {
	// TODO: Check against duplicate card fronts

	front, back := commandAddArgs(cmd)

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

	_, err = cmd.Exec("insert into schedulings (card_id, new, schedule_time, update_efactor) values (?, 1, datetime(), 1);", cardId)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to schedule card: %s.\n", err.Error())
		cmd.Exit(1)
	}

	fmt.Printf("%d\n", cardId)

	cmd.Commit()
}

