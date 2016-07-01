package main

import (
	"ivartj/args"
	"fmt"
	"io"
)

func init() {
	cmdRegister("help", helpCmd)
}

func helpCmdUsage(w io.Writer) {
	fmt.Fprintf(w, "Usage: %s help [ <subcommand> ]\n", mainProgramName)

	fmt.Fprintf(w, `
Description:
  Presents help messages for the given subcommands.
`)

	fmt.Fprintln(w, `
Options:
  -h, --help  Prints help message.
`)
}

func helpCmdArgs(cmd *cmdContext) (commands []string) {

	tok := args.NewTokenizer(cmd.Args)
	commands = []string{}

	for tok.Next() {

		if tok.IsOption() {

			switch tok.Arg() {

			case "-h", "--help":
				helpCmdUsage(cmd.Stdout)
				cmd.Exit(0)

			default:
				cmd.Fatalf("Unrecognized command, '%s'.\n", tok.Arg())
			}

		} else {
			commands = append(commands, tok.Arg())
		}

	}

	if tok.Err() != nil {
		cmd.Fatalf("Error occurred on processing command-line options: %s.\n", tok.Err().Error())
	}

	return commands
}

func helpCmd(cmd *cmdContext) {

	commands := helpCmdArgs(cmd)

	if len(commands) != 0 {
		subcmd := cmdNewContext(cmd.db)
		subcmd.Stdout = cmd.Stdout
		subcmd.Stdin = cmd.Stdin
		subcmd.Stderr = cmd.Stderr
		status := 0
		for _, command := range commands {
			status |= subcmd.Run(command, "-h")
		}
		cmd.Exit(status)
	} else {
		mainUsage(cmd.Stdout)
	}
}

