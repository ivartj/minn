package main

import (
	"testing"
)

func TestCmdMigrate(t *testing.T) {

	cmd := cmdNewContext(dbOpenTemp())

	status := cmd.Run("migrate", "-b", "-t", "0.1.1")
	if status != 0 {
		t.Fatalf("Failed to create database with initial schema version.\n")
	}

	status = cmd.Run("migrate")
	if status != 0 {
		t.Fatalf("Failed to migrate from initial schema version to current version.\n")
	}

}

