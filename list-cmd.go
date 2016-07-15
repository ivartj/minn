package main

import (
	"io"
	"fmt"
	"github.com/ivartj/minn/args"
)

func init() {
	cmdRegister("list", listCmd)
}

func listCmdUsage(w io.Writer) {
	fmt.Fprintf(w, "Usage: %s list\n", mainProgramName)

	fmt.Fprintln(w, `
Description:
  Lists cards by card IDs, ordered so that the earliest scheduled come first.
  Options can be used to control the selection of cards listed.

Options:
  -h, --help          Prints help message.
  --select-scheduled  Limit selection to scheduled cards.
  --select-new        Limit selection to new (unreviewed) cards.
`)

}

func listCmdArgs(cmd *cmdContext) (whereClause string) {

	tok := args.NewTokenizer(cmd.Args)
	whereExpressions := []string{}

	for tok.Next() {

		if !tok.IsOption() {
			cmd.Fatalf("Unexpected non-option, '%s'.\n", tok.Arg())
		}

		switch tok.Arg() {

		case "-h", "--help":
			listCmdUsage(cmd.Stdout)
			cmd.Exit(0)

		case "--select-scheduled":
			whereExpressions = append(whereExpressions, fmt.Sprintf("schedule_time < '%s'", cmd.Now().Format(utilTimeFormat)))

		case "--select-new":
			whereExpressions = append(whereExpressions, "state is 0")

		default:
			cmd.Fatalf("Unrecognized option, '%s'.", tok.Arg())

		}

	}

	if tok.Err() != nil {
		cmd.Fatalf("Error on processing command-line options: %s.\n", tok.Err().Error())
	}

	if len(whereExpressions) == 0 {
		whereClause = ""
	} else {
		whereClause = "where "
		for i, expr := range whereExpressions {
			if i != 0 {
				whereClause += " and "
			}
			whereClause += expr
		}
	}

	return whereClause
}

func listCmd(cmd *cmdContext) {
	whereClause := listCmdArgs(cmd)
	rows, err := cmd.Query("select card_id from cards " + whereClause + " order by schedule_time asc;")
	if err != nil {
		cmd.Fatalf("Database query failed: %s.\n", err.Error())
	}

	var cardID int
	for rows.Next() {
		err = rows.Scan(&cardID)
		if err != nil {
			cmd.Fatalf("Failed to scan card ID: %s.\n", err.Error())
		}
		cmd.Println(cardID)
	}
	cmd.Commit()
}

