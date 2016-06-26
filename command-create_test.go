package main

import (
	"testing"
)

func TestCommandCreate(t *testing.T) {

	cmd := cmdNewContext(dbOpenTemp())
	status := cmd.Run([]string{ "create" })
	if status != 0 {
		t.Fatalf("Status of command non-zero.")
	}

}

