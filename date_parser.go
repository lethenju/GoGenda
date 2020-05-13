/*
MIT License

Copyright (c) 2020 Julien LE THENO

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
/*
 ============= GOGENDA SOURCE CODE ===========
 @Description : GoGenda is a CLI for google agenda, to focus on one task at a time and logs your activity
 @Author : Julien LE THENO
 =============================================
*/
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
		return time.Date(time.Now().Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local), nil
	}
	date, err = time.Parse("2006/01/02", dateToParse)
	if err == nil {
		return date, nil
	}
	date, err = time.Parse("01/02", dateToParse)
	if err == nil {
		return time.Date(time.Now().Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local), nil
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

func buildDateFromDateTime(dateStr string, timeStr string, date *time.Time) (errTime error, errDate error) {
	*date, errTime = timeParser(timeStr)
	if errTime == nil { // we have our time
		t, errDate := dateParser(dateStr)
		if errDate == nil { // Date is correct
			*date = time.Date(t.Year(), t.Month(), t.Day(), date.Hour(), date.Minute(), date.Second(), 0, time.Local)
		}
	}
	return errTime, errDate
}

func buildDateFromTimeDate(timeStr string, dateStr string, date *time.Time) (errTime error, errDate error) {
	*date, errDate = dateParser(dateStr)
	if errDate == nil { // Date is correct
		t, errTime := timeParser(timeStr)
		if errTime == nil { // we have our time
			*date = time.Date(date.Year(), date.Month(), date.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.Local)
		}
	}
	return errTime, errDate
}
