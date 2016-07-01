package main

import (
	"testing"
)

func TestCmdRate(t *testing.T) {

	db := dbOpenTemp()
	cmd := cmdNewContext(db)

	status := cmd.Run("create")
	if status != 0 {
		t.Fatalf("Failed to initialize test database.\n")
	}

	status = cmd.Run("add", "FRONT1", "BACK2")
	if status != 0 {
		t.Fatalf("Failed to add test card.\n")
	}

	status = cmd.Run("add", "FRONT2", "BACK2")
	if status != 0 {
		t.Fatalf("Failed to add test card.\n")
	}

	status = cmd.Run("add", "FRONT3", "BACK3")
	if status != 0 {
		t.Fatalf("Failed to add test card.\n")
	}

	status = cmd.Run("rate", "5")
	if status != 0 {
		t.Fatalf("'rate 5' command failed.\n")
	}
	
	status = cmd.Run("rate", "3")
	if status != 0 {
		t.Fatalf("'rate 3' command failed.\n")
	}

	status = cmd.Run("rate", "0")
	if status != 0 {
		t.Fatalf("'rate 0' command failed.\n")
	}

	status = cmd.Run("rate", "5")
	if status != 0 {
		t.Fatalf("'rate 5' command failed.\n")
	}

	status = cmd.Run("rate", "5")
	if status != 0 {
		t.Fatalf("'rate 5' command failed.\n")
	}

	status = cmd.Run("current")
	if status == 0 {
		t.Fatalf("'current' command unexpectedly succeeded when there should be no more cards to rate.\n")
	}
}

