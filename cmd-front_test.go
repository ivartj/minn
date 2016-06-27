package main

import (
	"testing"
)

func TestCmdFront(t *testing.T) {

	db := dbOpenTemp()

	cmd := cmdNewContext(db)

	status := cmd.Run("create")
	if status != 0 {
		t.Fatalf("Fail to create test database.")
	}

	status = cmd.Run("add", "FRONT", "BACK")
	if status != 0 {
		t.Fatalf("Failed to add test card.\n")
	}

	status = cmd.Run("front")
	if status != 0 {
		t.Fatalf("Non-zero exit status on \"front\" command.")
	}

	status = cmd.Run("front", "1")
	if status != 0 {
		t.Fatalf("Non-zero exit status on \"front 1\" command.")
	}

}

