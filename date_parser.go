package main

import (
	"errors"
	"strings"
	"time"
)

// Parses a date given in parameters
// accepted input :
// case not sensitive : "yesterday", "today", "tommorow"
// date in YYYY-MM-DD, YYYY/MM/DD, MM-DD, MM/DD
func dateParser(dateToParse string) (date time.Time, err error) {
	date = time.Now()
	if strings.ToUpper(dateToParse) == "TODAY" {
		return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local), nil
	} else if strings.ToUpper(dateToParse) == "YESTERDAY" {
		date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
		return date.Add(-24 * time.Hour), nil
	} else if strings.ToUpper(dateToParse) == "TOMORROW" {
		date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
		return date.Add(24 * time.Hour), nil
	}
	date, err = time.Parse("2006-01-02", dateToParse)
	if err == nil {
		return date, nil
	}
	// Just the day
	date, err = time.Parse("02", dateToParse)
	if err == nil {
		return date, nil
	}
	date, err = time.Parse("01-02", dateToParse)
	if err == nil {
		return date, nil
	}
	date, err = time.Parse("2006/01/02", dateToParse)
	if err == nil {
		return date, nil
	}
	date, err = time.Parse("01/02", dateToParse)
	if err == nil {
		return date, nil
	}
	return time.Now(), errors.New("Wrong formatting")
}

func timeParser(timeStr string) (t time.Time, err error) {
	if strings.ToUpper(timeStr) == "NOW" {
		return time.Now(), nil
	}
	// Just the hour
	t, err = time.Parse("15", timeStr)
	if err == nil {
		return t, err
	}
	// Hour and minute
	t, err = time.Parse("15:04", timeStr)
	if err == nil {
		return t, err
	}
	// Hour minute and second
	t, err = time.Parse("15:04:05", timeStr)
	if err == nil {
		return t, err
	}
	return time.Now(), errors.New("Wrong formatting")
}
