package main

import (
	"fmt"
	"io"
	"ivartj/args"
)

func init() {
	cmdRegister("current", currentCmd)
}

func currentCmdUsage(w io.Writer) {
	fmt.Fprintf(w, "Usage: %s current\n", mainProgramName)

	fmt.Fprintf(w, `
Description:
  Presents the card ID of the current card.
`)

	fmt.Fprintln(w, `
Options:
  -h, --help  Prints help message.
`)
}

func currentCmdArgs(cmd *cmdContext) {

	tok := args.NewTokenizer(cmd.Args)

	for tok.Next() {

		if tok.IsOption() {

			switch tok.Arg() {
			case "-h", "--help":
				currentCmdUsage(cmd.Stdout)
				cmd.Exit(0)
			default:
				cmd.Fatalf("Unrecognized option, '%s'.\n", tok.Arg())
			}
				
		} else {
			currentCmdUsage(cmd.Stderr)
			cmd.Exit(1)
		}
	}

	if tok.Err() != nil {
		cmd.Fatalf("Error occurred on processing command-line options: %s.\n", tok.Err().Error())
	}
}

func currentCmd(cmd *cmdContext) {

	currentCmdArgs(cmd)

	cardId, err := utilCurrentCard(cmd)
	if err != nil {
		cmd.Fatalf("Failed to get current card ID: %s.\n", err.Error())
	}

	cmd.Println(cardId)
	
	cmd.Commit()
}

