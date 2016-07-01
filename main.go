package main

import (
	"ivartj/args"
	"os"
	"fmt"
	"io"
	"bytes"
	"bufio"
	"strconv"
)

const (
	mainProgramName		= "repete"
	mainProgramVersion	= "0.1-SNAPSHOT"

	mainSchemaVersion	= "0.1.2"
)

var (
	mainConfDeckPath		= mainProgramName + ".deck"
	mainConfMaxRelearnBacklog	= 10
)

func mainUsage(out io.Writer) {
	fmt.Fprintf(out, "Usage: %s [ <options> ] <subcommand> ...\n", mainProgramName)

	// Print the usage string of each subcommand
	cmd := cmdNewContext(dbOpenTemp())
	buf := bytes.NewBuffer([]byte{})
	cmd.Stdout = buf
	for _, command := range cmdList {
		buf.Reset()
		cmd.Run(command.name, "-h")
		line, _ := bufio.NewReader(buf).ReadString('\n')
		fmt.Fprint(out, line)
	}

	fmt.Fprintf(out, `
Description:
  %s is a command-line application which manages a deck of digital flash cards.
`, mainProgramName)

	fmt.Fprintln(out, `
Options before subcommand:
  -h, --help  Prints help message.
  --version   Prints version.
  -d, --deck  Specifies path of deck database to operate upon.
  -b, --max-relearn-backlog=<number>
	      Limits how far back in the backlog the current card can be.
`)

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

			case "-b", "--max-relearn-backlog":
				param, err := tok.TakeParameter()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Failed to get parameter to '%s' option: %s.\n", tok.Arg(), err.Error())
					os.Exit(1)
				}
				i, err := strconv.Atoi(param)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Failed to parse parameter to '%s' option: %s.\n", tok.Arg(), err.Error())
					os.Exit(1)
				}
				if i <= 0 {
					fmt.Fprintf(os.Stderr, "Parameter to '%s' option must be above zero.\n", tok.Arg())
					os.Exit(1)
				}
				mainConfMaxRelearnBacklog = i

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

