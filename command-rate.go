package main

import (
	"fmt"
	"os"
	"strconv"
	"io"
)

func commandRateUsage(w io.Writer) {
	fmt.Fprintf(w, "Usage: %s <deck> rate <rating>\n", mainProgramName)
}

func commandRateArgs(cmd *cmdContext) int {

	plainArgs := []string{}

	for cmd.Args.Next() {

		if cmd.Args.IsOption() {

			switch cmd.Args.Arg() {
			case "-h", "--help":
				commandRateUsage(os.Stdout)
				cmd.Exit(0)
			default:
				fmt.Fprintf(os.Stderr, "Unrecognized option, '%s'.\n", cmd.Args.Arg())
				cmd.Exit(1)
			}

		} else {
			plainArgs = append(plainArgs, cmd.Args.Arg())
		}
	}

	if cmd.Args.Err() != nil {
		fmt.Fprintf(os.Stderr, "Error on parsing command-line arguments: %s.\n", cmd.Args.Err().Error())
		cmd.Exit(1)
	}

	if len(plainArgs) != 1 {
		commandRateUsage(os.Stderr)
		cmd.Exit(1)
	}

	rating, err := strconv.Atoi(plainArgs[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse rating: %s.\n", err.Error())
		cmd.Exit(1)
	}

	if rating < 0 || rating > 5 {
		fmt.Fprintf(os.Stderr, "Rating must be 0 to 5.\n")
		cmd.Exit(1)
	}

	return rating
}

func commandRate(cmd *cmdContext) {

	rating := commandRateArgs(cmd)

	cardId, err := sm2CurrentCard(cmd.DB())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get current card ID: %s.\n", err.Error())
		cmd.Exit(1)
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
		fmt.Fprintf(os.Stderr, "Database error: %s.\n", err.Error())
		cmd.Exit(1)
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
			fmt.Fprintf(os.Stderr, "Failed to update E-factor: %s.\n", err.Error())
			cmd.Exit(1)
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
				fmt.Fprintf(os.Stderr, "Failed to update interval: %s.\n", err.Error())
				cmd.Exit(1)
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
		fmt.Fprintf(os.Stderr, "Failed to add new scheduling: %s.\n", err.Error())
		cmd.Exit(1)
	}

	cmd.Commit()
}

