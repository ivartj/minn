package main

import (
	"fmt"
	"io"
	"ivartj/args"
)

func init() {
	cmdRegister("add", addCmd)
}

func addCmdUsage(w io.Writer) {
	fmt.Fprintf(w, "Usage: %s add <front> <back>\n", mainProgramName)
}

func addCmdArgs(cmd *cmdContext) (string, string) {

	plainArgs := []string{}
	tok := args.NewTokenizer(cmd.Args)

	for tok.Next() {

		if tok.IsOption() {
			switch tok.Arg() {
			case "-h", "--help":
				addCmdUsage(cmd.Stdout)
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
		addCmdUsage(cmd.Stderr)
		cmd.Exit(1)
	}

	return plainArgs[0], plainArgs[1]
}

func addCmd(cmd *cmdContext) {
	// TODO: Check against duplicate card fronts

	front, back := addCmdArgs(cmd)

	nowStr := cmd.Now().Format(utilTimeFormat)
	res, err := cmd.Exec(`
		insert into
			cards (entry_time, front, back, efactor, interval, state, schedule_time)
			values (?, ?, ?, 2.5, 1, 0, ?);
	`, nowStr, front, back, nowStr)
	if err != nil {
		cmd.Fatalf("Database error: %s.\n", err.Error())
	}

	cardId, err := res.LastInsertId()
	if err != nil {
		cmd.Fatalf("Failed to get card ID from database engine: %s.\n", err.Error())
	}

	cmd.Println(cardId)

	cmd.Commit()
}

