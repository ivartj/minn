package main

import (
	"database/sql"
	"fmt"
	"errors"
)

func sm2CurrentCard(db *sql.Tx) (int, error) {

	row := db.QueryRow(`
		select card_id
			from schedulings
			where schedule_time <= datetime()
			order by schedule_time asc
			limit 1;`)
	cardId := 0
	err := row.Scan(&cardId)
	if err == sql.ErrNoRows {
		return 0, errors.New("No scheduled cards")
	}
	if err != nil {
		return 0, fmt.Errorf("Database error: %s.\n", err.Error())
	}

	return cardId, nil
}

