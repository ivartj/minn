package main

import (
	"testing"
	"bytes"
	"fmt"
)

func TestCmdRemove(t *testing.T) {

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
		t.Fatalf("Failed to read test card ID: %s.\n", err.Error())
	}

	status = cmd.Run("remove", cardId)
	if status != 0 {
		t.Fatalf("Non-zero exit status on removing card.\n")
	}

	status = cmd.Run("front", cardId)
	if status == 0 {
		t.Fatalf("'front' command unexpectedly succeeds on removed card ID.\n")
	}
}

