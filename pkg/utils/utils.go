package utils

import "time"

func ParseDate(dateStr *string) (time.Time, error) {
	if *dateStr == "" {
		*dateStr = time.Now().Format(time.DateOnly)
	}

	date, err := time.Parse(time.DateOnly, *dateStr)
	if err != nil {
		return time.Time{}, err
	}

	return date, nil
}

func EndOfDay(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 0, date.Location())
}
