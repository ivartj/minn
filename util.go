package main

import (
	"database/sql"
	"errors"
	"fmt"
)

const (
	utilTimeFormat = "2006-01-02 15:04:05"
)

type utilCardStatus int

const (
	CARD_NEW utilCardStatus			= 0
	CARD_RELEARN utilCardStatus		= 1
	CARD_REVIEW utilCardStatus		= 2
)

func utilCurrentCard(cmd *cmdContext) (int, error) {

	row := cmd.QueryRow(`
		select card_id
			from cards
			where schedule_time <= ?
			order by schedule_time asc
			limit 1;`, cmd.Now().Format(utilTimeFormat))
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

