package main

import (
	"fmt"
	"strconv"
	"io"
	"github.com/ivartj/minn/args"
	"time"
)

func init() {
	cmdRegister("rate", rateCmd)
}

func rateCmdUsage(w io.Writer) {
	fmt.Fprintf(w, "Usage: %s rate <rating>\n", mainProgramName)

	fmt.Fprintf(w, `
Description:
  Rates and reschedules the current card.

  The rating is a numeric which is 0 to 5 inclusive. The different numbers have
  the following meanings, taken from the specification of the SuperMemo 2
  algorithm:

  5) Perfect response.

  4) Correct response after a hesitation.

  3) Correct response recalled with serious difficulty.

  2) Incorrect response; where the correct one seemed easy to recall.

  1) Incorrect response; the correct one remembered.

  0) Complete blackout.
`)

	fmt.Fprintln(w, `
Options:
  -h, --help  Prints help message.
`)

}

func rateCmdArgs(cmd *cmdContext) int {

	plainArgs := []string{}
	tok := args.NewTokenizer(cmd.Args)

	for tok.Next() {

		if tok.IsOption() {

			switch tok.Arg() {
			case "-h", "--help":
				rateCmdUsage(cmd.Stdout)
				cmd.Exit(0)
			default:
				cmd.Fatalf("Unrecognized option, '%s'.\n", tok.Arg())
			}

		} else {
			plainArgs = append(plainArgs, tok.Arg())
		}
	}

	if tok.Err() != nil {
		cmd.Fatalf("Error on parsing command-line arguments: %s.\n", tok.Err().Error())
	}

	if len(plainArgs) != 1 {
		rateCmdUsage(cmd.Stderr)
		cmd.Exit(1)
	}

	rating, err := strconv.Atoi(plainArgs[0])
	if err != nil {
		cmd.Fatalf("Failed to parse rating: %s.\n", err.Error())
	}

	if rating < 0 || rating > 5 {
		cmd.Fatalf("Rating must be 0 to 5.\n")
	}

	return rating
}

func rateCmd(cmd *cmdContext) {

	rating := rateCmdArgs(cmd)

	cardId, err := utilCurrentCard(cmd)
	if err != nil {
		cmd.Fatalf("Failed to get current card ID: %s.\n", err.Error())
	}

	var (
		efactor float64
		interval int
		state utilCardStatus
		scheduleTime time.Time
	)
	row := cmd.QueryRow("select efactor, interval, state from cards where card_id = ?", cardId)
	err = row.Scan(&efactor, &interval, &state)
	if err != nil {
		cmd.Fatalf("Database error: %s.\n", err.Error())
	}

	switch state {
	case CARD_NEW: fallthrough;
	case CARD_RELEARN:
		if rating >= 4 {
			state = CARD_REVIEW
			scheduleTime = cmd.Now().Add(time.Hour * time.Duration(24 * interval))
		} else {
			scheduleTime = cmd.Now()
		}

	case CARD_REVIEW:
		efactor = efactor - 0.8 + 0.28 * float64(rating) - 0.02 * float64(rating) * float64(rating)
		if efactor < 1.3 {
			efactor = 1.3
		}

		if rating < 3 {
			interval = 1
		} else {
			interval = int(float64(interval) * efactor + 0.5)
		}

		if rating < 4 {
			state = CARD_RELEARN
			scheduleTime = cmd.Now()
		} else {
			scheduleTime = cmd.Now().Add(time.Hour * time.Duration(24 * interval))
		}
	}

	_, err = cmd.Exec(`
		update cards
		set
			efactor = ?,
			interval = ?,
			state = ?,
			schedule_time = ?
		where card_id = ?;`, efactor, interval, state, scheduleTime.Format(utilTimeFormat), cardId)

	if err != nil {
		cmd.Fatalf("Error occurred on updating card: %s.\n", err.Error())
	}

	cmd.Commit()
}

