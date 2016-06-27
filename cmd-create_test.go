package main

import (
	"testing"
)

func TestCmdCreate(t *testing.T) {

	cmd := cmdNewContext(dbOpenTemp())
	status := cmd.Run("create")
	if status != 0 {
		t.Fatalf("Status of command non-zero.")
	}

}

