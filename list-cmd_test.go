package main

import (
	"testing"
)

func TestCmdList(t *testing.T) {

	db := dbOpenTemp()

	cmd := cmdNewContext(db)
	status := cmd.Run("create")
	if status != 0 {
		t.Fatalf("Failed to initialize test database.")
	}

	status = cmd.Run("list", "--select-scheduled", "--select-new")
	if status != 0 {
		t.Fatalf("List command failed.")
	}
}

