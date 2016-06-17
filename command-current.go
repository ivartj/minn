package main

import (
	"fmt"
	"os"
	"io"
)

func commandCurrentUsage(w io.Writer) {
	fmt.Fprintf(w, "Usage: %s <deck> current\n", mainProgramName)
}

func commandCurrentArgs(cmd *cmdContext) {

	for cmd.Args.Next() {

		if cmd.Args.IsOption() {

			switch cmd.Args.Arg() {
			case "-h", "--help":
				commandCurrentUsage(os.Stdout)
				cmd.Exit(0)
			default:
				fmt.Fprintf(os.Stderr, "Unrecognized option, '%s'.\n", cmd.Args.Arg())
				cmd.Exit(1)
			}
				
		} else {
			commandCurrentUsage(os.Stderr)
			cmd.Exit(1)
		}
	}

	if cmd.Args.Err() != nil {
		fmt.Fprintf(os.Stderr, "Error occurred on processing command-line options: %s.\n", cmd.Args.Err().Error())
		cmd.Exit(1)
	}
}

func commandCurrent(cmd *cmdContext) {

	commandCurrentArgs(cmd)

	cardId, err := sm2CurrentCard(cmd.DB())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get current card ID: %s.\n", err.Error())
		cmd.Exit(1)
	}

	fmt.Println(cardId)
	
	cmd.Commit()
}

