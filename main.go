package main

import (
	"ivartj/args"
	"os"
	"fmt"
	"io"
	"bytes"
	"bufio"
)

const (
	mainProgramName		= "repete"
	mainProgramVersion	= "0.1-SNAPSHOT"

	mainSchemaVersion	= "0.1.1"

)

func mainUsage(out io.Writer) {
	fmt.Fprintf(out, "Usage: %s <deck> <command> ...\n", mainProgramName)

	// Print first line, the usage string, of each subcommand
	for _, v := range cmdList {
		buf := bytes.NewBuffer([]byte{})
		v.usage(buf)
		line, _ := bufio.NewReader(buf).ReadString('\n')
		fmt.Fprint(out, line)
	}
}

func mainArgs() (string, []string) {
	tok := args.NewTokenizer(os.Args)
	plainArgs := []string{}

	for tok.Next() {

		if tok.IsOption() {

			switch tok.Arg() {
			case "-h", "--help":
				mainUsage(os.Stdout)
				os.Exit(0)

			case "--version":
				fmt.Printf("%s version %s\n")
				os.Exit(0)

			default:
				fmt.Fprintf(os.Stderr, "Unrecognized option, %s.\n", tok.Arg())
				os.Exit(1)
			}

		} else {
			plainArgs = append(plainArgs, tok.Arg())
			if(len(plainArgs) == 2) {
				break
			}
		}

	}

	if tok.Err() != nil {
		fmt.Fprintf(os.Stderr, "Error occurred on processing command-line arguments: %s.\n", tok.Err().Error())
		os.Exit(1)
	}

	if(len(plainArgs) < 2) {
		mainUsage(os.Stderr)
		os.Exit(1)
	}

	argv, noMoreOptions := tok.Remainder()
	if noMoreOptions {
		fmt.Fprintf(os.Stderr, "No-more-options marker (--) not permitted before subcommand.\n")
		os.Exit(1)
	}

	return plainArgs[0], append([]string{plainArgs[1]}, argv...)
}

func main() {
	deckfilepath, argv := mainArgs()

	cmd := cmdNewContext(dbOpenFile(deckfilepath))
	iargv := make([]interface{}, len(argv))
	for i, v := range argv {
		iargv[i] = v
	}
	os.Exit(cmd.Run(iargv...))
}

