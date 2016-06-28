package main

import (
	"fmt"
	"io"
	"ivartj/args"
)

func init() {
	cmdRegister("current", cmdCurrent, cmdCurrentUsage)
}

func cmdCurrentUsage(w io.Writer) {
	fmt.Fprintf(w, "Usage: %s current\n", mainProgramName)
}

func cmdCurrentArgs(cmd *cmdContext) {

	tok := args.NewTokenizer(cmd.Args)

	for tok.Next() {

		if tok.IsOption() {

			switch tok.Arg() {
			case "-h", "--help":
				cmdCurrentUsage(cmd.Stdout)
				cmd.Exit(0)
			default:
				cmd.Fatalf("Unrecognized option, '%s'.\n", tok.Arg())
			}
				
		} else {
			cmdCurrentUsage(cmd.Stderr)
			cmd.Exit(1)
		}
	}

	if tok.Err() != nil {
		cmd.Fatalf("Error occurred on processing command-line options: %s.\n", tok.Err().Error())
	}
}

func cmdCurrent(cmd *cmdContext) {

	cmdCurrentArgs(cmd)

	cardId, err := utilCurrentCard(cmd)
	if err != nil {
		cmd.Fatalf("Failed to get current card ID: %s.\n", err.Error())
	}

	cmd.Println(cardId)
	
	cmd.Commit()
}

