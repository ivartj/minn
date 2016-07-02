package main

import (
	"github.com/ivartj/minn/args"
	"fmt"
	"io"
	"sort"
)

func init() {
	cmdRegister("list", listCmd)
}

func listCmdUsage(w io.Writer) {
	fmt.Fprintf(w, "Usage: %s list\n", mainProgramName)
}

func listCmdArgs(cmd *cmdContext) {

	tok := args.NewTokenizer(cmd.Args)

	for tok.Next() {

		if tok.IsOption() {

			switch tok.Arg() {
			case "-h", "--help":
				listCmdUsage(cmd.Stdout)
				cmd.Exit(0)
			default:
				cmd.Fatalf("Unrecognized option, '%s'.\n", tok.Arg())
			}

		} else {
			listCmdUsage(cmd.Stderr)
			cmd.Exit(1)
		}
	}

	if tok.Err() != nil {
		cmd.Fatalf("Error on processing command-line options: %s.\n", tok.Err().Error())
	}

}

func listCmd(cmd *cmdContext) {
	listCmdArgs(cmd)
	ls := make([]string, 0, len(cmd.DeclMap))
	for id, _ := range cmd.DeclMap {
		ls = ls[:len(ls)+1]
		ls[len(ls)-1] = id
	}
	sort.Strings(ls)
	for _, v := range ls {
		cmd.Println(v)
	}
}

