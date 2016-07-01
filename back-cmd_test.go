package main

import (
	"testing"
	"bytes"
	"fmt"
)

func TestCmdBack(t *testing.T) {

	db := dbOpenTemp()
	cmd := cmdNewContext(db)

	status := cmd.Run("create")
	if status != 0 {
		t.Fatalf("Failed to initialize test database.\n")
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
		t.Fatalf("Failed to scan card ID from add output.\n")
	}

	buf = bytes.NewBuffer([]byte{})
	cmd.Stdout = buf
	status = cmd.Run("back", cardId)
	if status != 0 {
		t.Fatalf("back command had non-zero exit status.\n")
	}

	if string(buf.Bytes()) != "BACK\n" {
		t.Fatalf("Output '%s' not what expected.\n", string(buf.Bytes()))
	}
}

