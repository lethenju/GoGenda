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
package gogendalib

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lethenju/gogenda/internal/configuration"
	"github.com/lethenju/gogenda/internal/utilities"
	"github.com/lethenju/gogenda/pkg/colors"
	api "github.com/lethenju/gogenda/pkg/google_agenda_api"
	"google.golang.org/api/calendar/v3"
)

func planCommand(command Command, srv *calendar.Service) (err error) {

	// command[1] == action
	// action could be SHOW, MOVE, DELETE, RENAME

	// Small helper function to check if the string is a possible action
	containActionFunc := func(str string) bool {
		actions := [4]string{"SHOW", "MOVE", "DELETE", "RENAME"}
		for _, a := range actions {
			if a == str {
				return true
			}
		}
		return false
	}
	// by default, the action is SHOW
	action := "SHOW"

	if len(command) > 1 {
		// if there is at least one argument (command[0] is PLAN anyway)
		if containActionFunc(strings.ToUpper(command[1])) {
			action = strings.ToUpper(command[1])
			command = command[1:]
		}
	}

	if action == "SHOW" {
		// init the plan structure
		var planBuffer utilities.Plan

		// Get plan of all day
		begin := time.Now()
		begin = time.Date(begin.Year(), begin.Month(), begin.Day(), 0, 0, 0, 0, time.Local)
		if len(command) > 1 {
			begin, err = utilities.DateParser(command[1])
			if err != nil {
				return err
			}

		}

		nbDays := 1
		if len(command) > 2 {
			// Number of days to do
			nbDays, err = strconv.Atoi(command[2])
			if err != nil {
				return errors.New("Wrong argument '" + command[2] + "', should be a number")
			}
		}
		end := begin.Add(time.Duration(24*nbDays) * time.Hour)

		cals, err := api.GetActivitiesBetweenDates(begin.Format(time.RFC3339), end.Format(time.RFC3339), srv)
		if cals == nil {
			colors.DisplayError("Error")
			return err
		}
		events := cals.Items

		var lastevent time.Time
		if len(events) > 0 {
			lastevent = time.Now()
		} else {
			colors.DisplayOk("No events found")
		}
		for i, event := range events {
			beginTime, _ := time.Parse(time.RFC3339, event.Start.DateTime)
			endTime, _ := time.Parse(time.RFC3339, event.End.DateTime)
			if beginTime.Day() != lastevent.Day() {
				colors.DisplayInfoHeading(" Events of " + beginTime.Format("01/02"))
			}
			color, _ := api.GetColorNameFromColorID(event.ColorId)
			category := configuration.GetNameFromColor(color)
			if category == "default" {
				category = ""
			}
			category += "]"
			category = fmt.Sprintf("[%-6s", category)
			colors.DisplayOk("[" + strconv.Itoa(i) + "] [ " + beginTime.Format("15:04") + " -> " + endTime.Format("15:04") + " ] " + category + " : " + event.Summary)
			lastevent = beginTime

			// fill our data
			var eventStored utilities.EventStored
			eventStored.Name = event.Summary
			eventStored.CalendarID = event.Id
			planBuffer.Events = append(planBuffer.Events, eventStored)
		}
		// store our data
		utilities.StorePlan(&planBuffer)
		return err
	}
	// Load the plan data
	planBuffer, err := utilities.LoadPlan()

	if err != nil {
		return errors.New("Please call 'PLAN SHOW' first before modifying an event we dont know about")
	}
	// Now we need to get the id
	if len(command) == 1 {
		return errors.New("Please give an id of an event to modify")
	}
	index, err := strconv.Atoi(command[1])
	if err != nil {

		return errors.New("Id should be a number, you gave :" + command[1])
	}

	switch action {
	case "MOVE":
		// grab the old date and time
		date, err := api.GetStartDateForEventID(planBuffer.Events[index].CalendarID, srv)
		var t time.Time
		// parse date and time
		if len(command) == 3 || len(command) == 4 {
			// PLAN ID date
			// We need to change the date but not the time
			dateParsed, err := utilities.DateParser(command[2])
			date = time.Date(dateParsed.Year(), dateParsed.Month(), dateParsed.Day(), date.Hour(), date.Minute(), date.Second(), 0, time.Local)
			if len(command) == 4 {
				// PLAN ID date time
				t, err = utilities.TimeParser(command[3])
			}
			if err != nil {
				// PLAN ID time
				t, err = utilities.TimeParser(command[2])
				if len(command) == 4 {
					// PLAN ID time date
					dateParsed, err = utilities.DateParser(command[3])
					date = time.Date(dateParsed.Year(), dateParsed.Month(), dateParsed.Day(), date.Hour(), date.Minute(), date.Second(), 0, time.Local)

				}
				if err != nil {
					// incorrect date or time
					return err
				}
			}
		}
		if !t.IsZero() {
			date = time.Date(date.Year(), date.Month(), date.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.Local)
		}

		colors.DisplayOk("Moving element nb " + strconv.Itoa(index) + " : " + planBuffer.Events[index].Name + "to date and time " + date.Format(time.UnixDate))
		// Todo ask user if okay
		colors.DisplayError("MOVE NOT IMPLEMENTED YET")
		return err
	case "DELETE":
		colors.DisplayOk("Removing element nb " + strconv.Itoa(index) + " : " + planBuffer.Events[index].Name)
		// Todo ask user if okay
		colors.DisplayError("DELETE NOT IMPLEMENTED YET")
		return err
	case "RENAME":
		// Todo get the new name
		colors.DisplayOk("Renaming element nb " + strconv.Itoa(index) + " : " + planBuffer.Events[index].Name)
		// Todo ask user if okay
		colors.DisplayError("MOVE NOT IMPLEMENTED YET")
		return err
	}
	return err
}

// Add an event sometime
// If you want to add it now, you better use startCommand
func addCommand(command Command, srv *calendar.Service) (err error) {

	var date time.Time
	var endDate time.Time
	var name string
	var category string

	var isDateSet bool
	var isTimeSet bool
	var isEndDateSet bool

	askDate := func(date *time.Time) {
		askAgain := true
		for askAgain {
			inputStr := utilities.InputFromUser("date of event")
			t, err := utilities.DateParser(inputStr)
			if err != nil {
				colors.DisplayError("Wrong formatting !")
			} else {
				*date = time.Date(t.Year(), t.Month(), t.Day(), date.Hour(), date.Minute(), date.Second(), 0, time.Local)
				askAgain = false
			}
		}
	}
	askTime := func(date *time.Time) {
		askAgain := true
		for askAgain {
			inputStr := utilities.InputFromUser("begin time of event")
			t, err := utilities.TimeParser(inputStr)
			if err != nil {
				colors.DisplayError("Wrong formatting !")
			} else {
				*date = time.Date(date.Year(), date.Month(), date.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.Local)
				askAgain = false
			}
		}
	}
	askEndTime := func(endDate *time.Time) {
		askAgain := true
		for askAgain {
			inputStr := utilities.InputFromUser("end time of event")
			t, err := utilities.TimeParser(inputStr)
			if err != nil {
				colors.DisplayError("Wrong formatting !")
			} else {
				*endDate = time.Date(date.Year(), date.Month(), date.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.Local)
				if !endDate.After(date) {
					colors.DisplayError("End time cannot be before start time !")
				} else {
					askAgain = false
				}
			}
		}
	}
	askName := func(name *string) {
		*name = utilities.InputFromUser("name of event")
	}
	askCategory := func(category *string) {
		*category = utilities.InputFromUser("category of event")
	}

	if len(command) == 2 {
		// One argument given, we need to check which one is it : time, date or name
		date, err = utilities.TimeParser(command[1])
		if err == nil { // we have our time
			isTimeSet = true // So we set this flag on
		} else {
			date, err = utilities.DateParser(command[1])
			if err == nil { // We have our date
				isDateSet = true // So we set this flag on
			} else { // Its a category
				category = command[1]
			}
		}
	} else if len(command) == 3 {
		// Two arguments given, we need to check if they are :

		// time date
		// time category
		// date time
		// date category

		errTime, errDate := utilities.BuildDateFromTimeDate(command[1], command[2], &date)
		if errTime == nil && errDate == nil {
			isTimeSet = true
			isDateSet = true
		} else if errTime == nil && errDate != nil { // getting date failed
			isTimeSet = true
			category = command[2]
		} else {
			// date time
			// date category
			errTime, errDate := utilities.BuildDateFromDateTime(command[1], command[2], &date)
			if errTime == nil && errDate == nil {
				isTimeSet = true
				isDateSet = true
			} else if errDate == nil && errTime != nil { // getting time failed
				isDateSet = true
				category = command[2]
			}
		}

	} else if len(command) >= 4 {
		// Three arguments given. Need to verify if they are :

		// time date endDate
		// time date category
		// time category name
		// date time endDate
		// date time category
		// date category name

		errTime, errDate := utilities.BuildDateFromTimeDate(command[1], command[2], &date)
		if errTime == nil && errDate == nil {
			// time date category
			isTimeSet = true
			isDateSet = true

			endDate, err = utilities.TimeParser(command[3])
			if err == nil { // we have the end hour. Still need to fix the day
				endDate = time.Date(date.Year(), date.Month(), date.Day(), endDate.Hour(), endDate.Minute(), endDate.Second(), 0, time.Local)
				if date.After(endDate) { // like if the user wanted an event between 2 days (23:00 -> 01:00)
					endDate.Add(24 * time.Hour) // move to the next day
				}
				isEndDateSet = true // So we set this flag on
			} else {
				category = command[3]
			}
		} else if errTime == nil && errDate != nil { // getting date failed
			// time category name
			isTimeSet = true
			category = command[2]
			name = command[3]
		} else {

			// check if date time endDate
			// check if date time category
			// check if date category name
			errTime, errDate := utilities.BuildDateFromDateTime(command[1], command[2], &date)
			if errTime == nil && errDate == nil {
				// date time category
				isTimeSet = true
				isDateSet = true
				// Check if date time endDate
				endDate, err = utilities.TimeParser(command[3])
				if err == nil { // we have the end hour. Still need to fix the day
					endDate = time.Date(date.Year(), date.Month(), date.Day(), endDate.Hour(), endDate.Minute(), endDate.Second(), 0, time.Local)
					if date.After(endDate) { // like if the user wanted an event between 2 days (23:00 -> 01:00)
						endDate.Add(24 * time.Hour) // move to the next day
					}
					isEndDateSet = true // So we set this flag on
				} else {
					category = command[3]
				}
			} else if errDate == nil && errTime != nil { // getting time failed
				// check if date category name
				isDateSet = true
				category = command[2]
				name = command[3]
			}
		}
	}
	if len(command) >= 5 {
		// Now we have at least 4 arguments. They can be either :
		// time date endDate category
		// time date category name
		// date time endDate category
		// date time category name

		// Just check the last one as we already did the other arguments

		if isTimeSet && isDateSet && isEndDateSet {
			category = command[4]
			if len(command) > 5 {
				name = strings.Join(command[5:], " ")
			}
		} else if isTimeSet && isDateSet && category != "" {
			name = strings.Join(command[5:], " ")
		}
	}

	if !isDateSet {
		askDate(&date)
		isDateSet = true
	}
	if !isTimeSet {
		askTime(&date)
		isTimeSet = true
	}
	if !isEndDateSet {
		askEndTime(&endDate)
		isEndDateSet = true
	}
	if category == "" {
		askCategory(&category)
	}
	if name == "" {
		askName(&name)
	}

	color := configuration.GetColorFromName(category)
	colors.DisplayOk("Adding event " + name + " of category " + category + " starting " + date.Format("2006-01-02") + " at " + date.Format("15:04") + " until " + endDate.Format("15:04"))
	_, err = api.InsertActivity(name, color, date, endDate, srv)
	if err != nil {
		colors.DisplayError(err.Error())
	}
	return err
}
