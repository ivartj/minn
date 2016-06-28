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

	mainSchemaVersion	= "0.1.2"
)

var (
	mainConfDeckPath	= mainProgramName + ".deck"
)

func mainUsage(out io.Writer) {
	fmt.Fprintf(out, "Usage: %s [ -d <deck> ] <command> ...\n", mainProgramName)

	// Print first line, the usage string, of each subcommand
	for _, v := range cmdList {
		buf := bytes.NewBuffer([]byte{})
		v.usage(buf)
		line, _ := bufio.NewReader(buf).ReadString('\n')
		fmt.Fprint(out, line)
	}
}

func mainArgs() ([]string) {
	tok := args.NewTokenizer(os.Args)
	var command string
	commandGiven := false

	for tok.Next() {

		if tok.IsOption() {

			switch tok.Arg() {
			case "-h", "--help":
				mainUsage(os.Stdout)
				os.Exit(0)

			case "--version":
				fmt.Printf("%s version %s\n")
				os.Exit(0)

			case "-d", "--deck":
				param, err := tok.TakeParameter()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Failed to get parameter to '%s' option: %s.\n", tok.Arg(), err.Error())
					os.Exit(1)
				}
				mainConfDeckPath = param

			default:
				fmt.Fprintf(os.Stderr, "Unrecognized option, %s.\n", tok.Arg())
				os.Exit(1)
			}

		} else {
			command = tok.Arg()
			commandGiven = true
			break
		}

	}

	if tok.Err() != nil {
		fmt.Fprintf(os.Stderr, "Error occurred on processing command-line arguments: %s.\n", tok.Err().Error())
		os.Exit(1)
	}

	if(!commandGiven) {
		mainUsage(os.Stderr)
		os.Exit(1)
	}

	argv, noMoreOptions := tok.Remainder()
	if noMoreOptions {
		fmt.Fprintf(os.Stderr, "No-more-options marker (--) not permitted before subcommand.\n")
		os.Exit(1)
	}

	return append([]string{command}, argv...)
}

func main() {
	argv := mainArgs()

	cmd := cmdNewContext(dbOpenFile(mainConfDeckPath))
	iargv := make([]interface{}, len(argv))
	for i, v := range argv {
		iargv[i] = v
	}
	os.Exit(cmd.Run(iargv...))
}

