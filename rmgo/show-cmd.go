package main

import (
	"github.com/ivartj/minn/args"
	"fmt"
	"go/printer"
	"io"
)

func init() {
	cmdRegister("show", showCmd)
}

func showCmdUsage(w io.Writer) {
	fmt.Fprintf(w, "Usage: %s show <identifier> ...\n", mainProgramName)
}

func showCmdArgs(cmd *cmdContext) (idents []string) {

	tok := args.NewTokenizer(cmd.Args)
	idents = []string{}

	for tok.Next() {

		if tok.IsOption() {

			switch tok.Arg() {
			case "-h", "--help":
				showCmdUsage(cmd.Stdout)
				cmd.Exit(0)

			default:
				cmd.Fatalf("Unrecognized option, '%s'.\n", tok.Arg())
			}

		} else {
			idents = append(idents, tok.Arg())
		}

	}

	if tok.Err() != nil {
		cmd.Fatalf("Error on processing command-line options: %s.\n", tok.Err().Error())
	}

	if len(idents) == 0 {
		showCmdUsage(cmd.Stderr)
		cmd.Exit(1)
	}

	return idents
}

func showCmd(cmd *cmdContext) {

	idents := showCmdArgs(cmd)

	for _, ident := range idents {
		decl, ok := cmd.DeclMap[ident]
		if !ok {
			cmd.Fatalf("Could not find declaration for '%s'.\n", ident)
		}
		err := printer.Fprint(cmd.Stdout, cmd.FileSet, decl)
		if err != nil {
			cmd.Fatalf("Error occurred on printing declaration: %s.\n", err.Error())
		}
		cmd.Println()
	}
}

