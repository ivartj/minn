package main

import (
	"testing"
)

func TestCmdAdd(t *testing.T) {

	db := dbOpenTemp()

	cmd := cmdNewContext(db)
	status := cmd.Run("create")
	if status != 0 {
		t.Fatalf("Failed to initialize test database.")
	}

	cmd = cmdNewContext(db)
	status = cmd.Run("add", "FRONT", "BACK")
	if status != 0 {
		t.Fatalf("Command status non-zero.")
	}

}

