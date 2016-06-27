package main

import (
	"fmt"
	"strconv"
	"io"
	"ivartj/args"
)

func init() {
	cmdRegister("rate", cmdRate, cmdRateUsage)
}

func cmdRateUsage(w io.Writer) {
	fmt.Fprintf(w, "Usage: %s rate <rating>\n", mainProgramName)
}

func cmdRateArgs(cmd *cmdContext) int {

	plainArgs := []string{}
	tok := args.NewTokenizer(cmd.Args)

	for tok.Next() {

		if tok.IsOption() {

			switch tok.Arg() {
			case "-h", "--help":
				cmdRateUsage(cmd.Stdout)
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
		cmdRateUsage(cmd.Stderr)
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

func cmdRate(cmd *cmdContext) {

	rating := cmdRateArgs(cmd)

	cardId, err := sm2CurrentCard(cmd.DB())
	if err != nil {
		cmd.Fatalf("Failed to get current card ID: %s.\n", err.Error())
	}

	var (
		efactor float64
		interval int
		updateEfactor bool
		updateInterval bool
	)
	row := cmd.QueryRow("select efactor, interval, update_efactor, update_interval from cards natural join schedulings where card_id = ?", cardId)
	err = row.Scan(&efactor, &interval, &updateEfactor, &updateInterval)
	if err != nil {
		cmd.Fatalf("Database error: %s.\n", err.Error())
	}

	if updateEfactor {
		efactor = efactor - 0.8 + 0.28 * float64(rating) - 0.02 * float64(rating) * float64(rating)
		if efactor < 1.3 {
			efactor = 1.3
		}
		_, err = cmd.Exec(`
			update cards
				set efactor = ?
				where card_id = ?;
		`, efactor, cardId)
		if err != nil {
			cmd.Fatalf("Failed to update E-factor: %s.\n", err.Error())
		}
	}

	if updateInterval {
		newInterval := 0
		if rating >= 3 {
			switch interval {
			case 0:
				newInterval = 1
			case 1:
				newInterval = 6
			default:
				newInterval = int(float64(interval) * efactor + 0.5) // rounding
			}
		} else {
			newInterval = 0
		}

		if newInterval != interval {
			_, err := cmd.Exec(`
				update cards
					set interval = ?
					where card_id = ?;
			`, newInterval, cardId)
			if err != nil {
				cmd.Fatalf("Failed to update interval: %s.\n", err.Error())
			}
		}
		interval = newInterval
	}

	// "After each repetition on a given day repeat again all items that
	// scored below four in the quality assessment"
	updateInterval = true
	if rating < 4 {
		interval = 0
		updateInterval = false
	}

	_, err = cmd.Exec(`
		insert or replace into
			schedulings (card_id, new, schedule_time, update_efactor, update_interval)
		values
			(?, 0, datetime('now', ?), ?, ?);
	`, cardId, fmt.Sprintf("+%d day", interval), rating >= 3, updateInterval)
	if err != nil {
		cmd.Fatalf("Failed to add new scheduling: %s.\n", err.Error())
		cmd.Exit(1)
	}

	cmd.Commit()
}

