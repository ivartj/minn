package main

import (
	"bytes"
	"strings"
	"fmt"
	"testing"
)

func TestHelpOption(t *testing.T) {

	cmd := cmdNewContext(dbOpenTemp())
	outbuf := bytes.NewBuffer([]byte{})
	errbuf := bytes.NewBuffer([]byte{})

	for _, command := range cmdList {

		outbuf.Reset()
		errbuf.Reset()

		cmd.Stdout = outbuf
		cmd.Stderr = errbuf
		
		for _, option := range []string{ "-h", "--help" } {

			status := cmd.Run(command.name, option)
			if status != 0 {
				t.Fatalf("Non-zero exit status when passing %s to %s subcommand.\n", option, command.name) 
			}

			if errbuf.Len() != 0 {
				t.Fatalf("Subcommand %s wrote to stderr when passed %s option.\n", command.name, option)
			}

			if !strings.HasPrefix(outbuf.String(), fmt.Sprintf("Usage: %s %s", mainProgramName, command.name)) {
				t.Fatalf("The output of %s %s did not have the expected initial output.\n", command.name, option)
			}

			if !strings.Contains(outbuf.String(), "\n") {
				t.Fatalf("The output of %s %s does not contain a newline.\n", command.name, option)
			}

			if !strings.HasSuffix(outbuf.String(), "\n") {
				t.Fatalf("The output of %s %s does not end of newline.\n", command.name, option)
			}
		}
	}
}

