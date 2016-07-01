package main

import (
	"testing"
	"bytes"
	"fmt"
	"os"
)

func TestCmdCurrent(t *testing.T) {

	db := dbOpenTemp()
	cmd := cmdNewContext(db)

	status := cmd.Run("create")
	if status != 0 {
		t.Fatalf("Failed to initialize test database.\n")
	}

	status = cmd.Run("add", "FRONT", "BACK")
	if status != 0 {
		t.Fatalf("Failed to add test card.\n")
	}

	buf := bytes.NewBuffer([]byte{})
	cmd.Stdout = buf
	status = cmd.Run("current")
	if status != 0 {
		t.Fatalf("Non-zero exit status.\n")
	}

	var cardId int
	_, err := fmt.Fscanln(buf, &cardId)
	if err != nil {
		t.Fatalf("Error occurred on scanning output from command: %s.\n", err.Error())
	}

	cmd.Stdout = os.Stdout
	status = cmd.Run("front", fmt.Sprint(cardId))
	if status != 0 {
		t.Fatalf("Failed to ensure existence of current card ID.\n")
	}
	
}

