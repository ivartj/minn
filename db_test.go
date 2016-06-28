package main

import (
	"testing"
	"fmt"
)

func TestDbGetMigrationPath(t *testing.T) {

	dbRegisterMigration("T.0.1", "T.0.2", "A")
	dbRegisterMigration("T.0.2", "T.0.3", "B")

	path, err := dbGetMigrationPath("T.0.1", "T.0.3")
	if err != nil {
		t.Fatalf("Could not discover registered test path: %s.\n", err.Error())
	}

	for _, m := range path {
		fmt.Printf("%s -> %s\n", m.from, m.to)
	}

}

