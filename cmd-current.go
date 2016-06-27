package main

import (
	"fmt"
	"os"
	"io"
	"ivartj/args"
)

func init() {
	cmdRegister("current", cmdCurrent, cmdCurrentUsage)
}

func cmdCurrentUsage(w io.Writer) {
	fmt.Fprintf(w, "Usage: %s <deck> current\n", mainProgramName)
}

func cmdCurrentArgs(cmd *cmdContext) {

	tok := args.NewTokenizer(cmd.Args)

	for tok.Next() {

		if tok.IsOption() {

			switch tok.Arg() {
			case "-h", "--help":
				cmdCurrentUsage(os.Stdout)
				cmd.Exit(0)
			default:
				fmt.Fprintf(os.Stderr, "Unrecognized option, '%s'.\n", tok.Arg())
				cmd.Exit(1)
			}
				
		} else {
			cmdCurrentUsage(os.Stderr)
			cmd.Exit(1)
		}
	}

	if tok.Err() != nil {
		fmt.Fprintf(os.Stderr, "Error occurred on processing command-line options: %s.\n", tok.Err().Error())
		cmd.Exit(1)
	}
}

func cmdCurrent(cmd *cmdContext) {

	cmdCurrentArgs(cmd)

	cardId, err := sm2CurrentCard(cmd.DB())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get current card ID: %s.\n", err.Error())
		cmd.Exit(1)
	}

	fmt.Println(cardId)
	
	cmd.Commit()
}

