package main

import (
	"testing"
	"bytes"
	"fmt"
)

func TestCmdFront(t *testing.T) {

	db := dbOpenTemp()

	cmd := cmdNewContext(db)

	status := cmd.Run("create")
	if status != 0 {
		t.Fatalf("Fail to create test database.")
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
		t.Fatalf("Failed to scan card ID from add command output.\n")
	}

	buf = bytes.NewBuffer([]byte{})
	cmd.Stdout = buf
	status = cmd.Run("front")
	if status != 0 {
		t.Fatalf("Non-zero exit status on \"front\" command.")
	}

	if string(buf.Bytes()) != "FRONT\n" {
		t.Fatalf("Output '%s' not expected of 'front' command.\n", string(buf.Bytes()))
	}

	buf = bytes.NewBuffer([]byte{})
	cmd.Stdout = buf
	status = cmd.Run("front", cardId)
	if status != 0 {
		t.Fatalf("Non-zero exit status on \"front 1\" command.")
	}

	if string(buf.Bytes()) != "FRONT\n" {
		t.Fatalf("Output '%s' not expected of 'front %d' command.\n", string(buf.Bytes()), cardId)
	}

}

