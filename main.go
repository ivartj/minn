package main

import (
	"github.com/ivartj/norske-irc-kanaler.com/args"
	"os"
	"fmt"
	"io"
)

const (
	mainProgramName		= "repete"
	mainProgramVersion	= "0.1-SNAPSHOT"

	mainSchemaVersion	= "0.1.1"

)

var (
	mainCommandMap		= map[string]func(*cmdContext){
		"create" : commandCreate,
		"current" : commandCurrent,
		"front" : commandFront,
		"back" : commandBack,
		"rate" : commandRate,
		"add" : commandAdd,
		"edit" : commandEdit,
		"remove" : commandRemove,
	}
)

func mainUsage(out io.Writer) {
	fmt.Fprintf(out, "Usage: %s <deck> <command> ...\n", mainProgramName)
	for _, v := range mainCommandMap {
		tok := args.NewTokenizer([]string{"command", "-h"})
		cmd := cmdNewContext("", tok)
		func() { defer func() { recover() }(); v(cmd) }()
	}
}

func mainArgs() (string, string, *args.Tokenizer) {
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

	return plainArgs[0], plainArgs[1], tok
}

func main() {
	deckfilepath, command, tok := mainArgs()

	cmdfn, ok := mainCommandMap[command]
	if !ok {
		fmt.Fprintf(os.Stderr, "Unrecognized command, %s.\n", command)
		os.Exit(1)
	} 

	cmd := cmdNewContext(deckfilepath, tok)

	status := 0
	defer func() {
		x := recover()
		switch x.(type) {
		case cmdExit:
			status = int(x.(cmdExit))
		case error:
			err := x.(error)
			fmt.Fprintf(os.Stderr, "Panic occurred: %s.\n", err.Error())
			status = 1
		}
		cmd.Rollback()
		os.Exit(status)
	}()
	cmdfn(cmd)
}

