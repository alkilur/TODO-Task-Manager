package http_server

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {

	if repeat == "" {
		return "", errors.New("'repeat' is empty")
	}

	nextDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", fmt.Errorf("invalid 'date' format: %w", err)
	}

	repeatParts := strings.Split(repeat, " ")
	switch repeatParts[0] {
	case "d":
		if len(repeatParts) != 2 {
			return "", fmt.Errorf("invalid 'repeat' format: %w", err)
		}
		days, err := strconv.Atoi(repeatParts[1])
		if err != nil || days >= 400 {
			return "", fmt.Errorf("invalid 'repeat' format: %w", err)
		}
		nextDate = nextDate.AddDate(0, 0, days)
		for nextDate.Before(now) || nextDate.Equal(now) {
			nextDate = nextDate.AddDate(0, 0, days)
		}
	case "y":
		nextDate = nextDate.AddDate(1, 0, 0)
		for nextDate.Before(now) || nextDate.Equal(now) {
			nextDate = nextDate.AddDate(1, 0, 0)
		}
	default:
		return "", fmt.Errorf("invalid 'repeat' format: %w", err)
	}

	return nextDate.Format("20060102"), nil
}
