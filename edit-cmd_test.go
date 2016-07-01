package main

import (
	"testing"
	"bytes"
	"fmt"
)

func TestCmdEdit(t *testing.T) {

	cmd := cmdNewContext(dbOpenTemp())
	status := cmd.Run("create")
	if status != 0 {
		t.Fatalf("Failed to create test database.\n")
	}

	buf := bytes.NewBuffer([]byte{})
	cmd.Stdout = buf
	status = cmd.Run("add", "FRONT", "BACK")
	if status != 0 {
		t.Fatalf("Failed to add test card.\n")
	}

	var cardId int
	_, err := fmt.Fscanln(buf, &cardId)
	if err != nil {
		t.Fatalf("Failed to scan added card ID.\n")
	}

	status = cmd.Run("edit", cardId, "NEW_FRONT", "NEW_BACK")
	if status != 0 {
		t.Fatalf("Failed to edit test card.\n")
	}

	buf.Reset()
	status = cmd.Run("front", cardId)
	if status != 0 {
		t.Fatalf("front command failed on edited test card.\n")
	}

	if buf.String() != "NEW_FRONT\n" {
		t.Fatalf("Output of front %d command does not appear to reflect edit.\n", cardId)
	}

	buf.Reset()
	status = cmd.Run("back", cardId)
	if status != 0 {
		t.Fatalf("back command failed on edited test card.\n")
	}

	if buf.String() != "NEW_BACK\n" {
		t.Fatalf("Output of back %d command does not appear to reflect edit.\n", cardId)
	}
}

