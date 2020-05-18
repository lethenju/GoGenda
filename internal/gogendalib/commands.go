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
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lethenju/gogenda/internal/configuration"
	"github.com/lethenju/gogenda/internal/current_activity"
	"github.com/lethenju/gogenda/pkg/colors"
	api "github.com/lethenju/gogenda/pkg/google_agenda_api"

	"google.golang.org/api/calendar/v3"
)

// Command : A command as a suite of arguments given by the user
type Command []string

// Add an event now
func startCommand(command Command, srv *calendar.Service) (err error) {
	var nameOfEvent string
	color := configuration.GetColorFromName(command[1])
	if len(command) == 2 && color != "blue" {
		fmt.Print(command)
		fmt.Print("Enter name of event :")
		scanner := bufio.NewScanner(os.Stdin)
		if !scanner.Scan() {
			return
		}
		nameOfEvent = scanner.Text()
		currentActivity, err := current_activity.GetCurrentActivity()
		if err == nil {
			// Stop the current activity
			err = api.StopActivity(currentActivity, srv)
			if err != nil {
				colors.DisplayError("There was an issue deleting the current event.")
			}
		}
	} else if len(command) == 2 {
		nameOfEvent = strings.Join(command[1:], " ")
	} else {
		nameOfEvent = strings.Join(command[2:], " ")
	}

	currentActivity, err := api.InsertActivity(nameOfEvent, color, time.Now(), time.Now().Add(30*time.Minute), srv)
	if err != nil {
		return err
	}
	current_activity.SetCurrentActivity(currentActivity)

	colors.DisplayOk("Successfully added activity ! ")
	return nil
}

func stopCommand(srv *calendar.Service) (err error) {

	currentActivity, err := current_activity.GetCurrentActivity()
	if err != nil {
		return errors.New("Nothing to stop")
	}

	duration, err := api.GetDuration(currentActivity)
	if err != nil {
		return err
	}
	colors.DisplayInfo("The activity '" + currentActivity.Summary + "' lasted " + duration)
	if err != nil {
		return err
	}
	err = api.StopActivity(currentActivity, srv)

	colors.DisplayOk("Successfully stopped the activity ! I hope it went well ")
	return nil
}

func deleteCommand(ctx *gogendaContext) (err error) {

	if ctx.activity.Id == "" {
		// Nothing to stop
		return errors.New("Nothing to delete")
	}
	err = deleteActivity(ctx.activity, ctx.srv)
	if err != nil {
		return err
	}
	displayOk(ctx, "Successfully deleted the activity ! ")
	return nil
}

func renameCommand(command Command, ctx *gogendaContext) (err error) {
	if ctx.activity.Id == "" {
		return errors.New("You dont have a current activity to rename")
	}
	var nameOfEvent string
	if len(command) == 1 {

		fmt.Println(ctx, "Enter name of event :  ")
		scanner := bufio.NewScanner(os.Stdin)
		if !scanner.Scan() {
			return
		}
		nameOfEvent = scanner.Text()
	} else {
		nameOfEvent = strings.Join(command[1:], " ")
	}
	err = renameActivity(ctx.activity, nameOfEvent, ctx.srv)
	if err != nil {
		return err
	}
	displayOk(ctx, "Successfully renamed the activity")
	return nil
}

func planCommand(command Command, ctx *gogendaContext) (err error) {
	// Get plan of all day
	begin := time.Now()
	begin = time.Date(begin.Year(), begin.Month(), begin.Day(), 0, 0, 0, 0, time.Local)
	if len(command) > 1 {
		begin, err = dateParser(command[1])
		if err != nil {
			return err
		}

	}

	end := time.Date(begin.Year(), begin.Month(), begin.Day(), 23, 59, 59, 0, time.Local)

	cals, err := getActivitiesBetweenDates(begin.Format(time.RFC3339), end.Format(time.RFC3339), ctx.srv)
	displayInfoHeading(ctx, " Events of "+begin.Format("01/02"))
	if cals == nil {
		displayError(ctx, "Error")
		return err
	}
	events := cals.Items
	for _, event := range events {
		beginTime, _ := time.Parse(time.RFC3339, event.Start.DateTime)
		endTime, _ := time.Parse(time.RFC3339, event.End.DateTime)
		color, _ := getColorNameFromColorID(event.ColorId)
		category := confGetNameFromColor(color, ctx.configuration)
		if category == "default" {
			category = ""
		}
		category += "]"
		category = fmt.Sprintf("[%-6s", category)
		displayOk(ctx, " [ "+beginTime.Format("15:04")+" -> "+endTime.Format("15:04")+" ] "+category+" : "+event.Summary)
	}

	if len(command) > 2 {
		// Number of days to do
		nbToRedo, err := strconv.Atoi(command[2])
		if err != nil {
			return err
		}
		// Removing a day (that has just be done)
		nbToRedo--
		// If we got to zero, we can return, its fine
		if nbToRedo == 0 {
			return nil
		}
		// adding a day
		begin = begin.Add(24 * time.Hour)
		// Putting it in the command
		command[1] = begin.Format("2006/01/02")
		command[2] = strconv.Itoa(nbToRedo)
		// Recursive call
		planCommand(command, ctx)
	}
	return err
}

// Add an event sometime
// If you want to add it now, you better use startCommand
func addCommand(command Command, ctx *gogendaContext) (err error) {

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
			inputStr := inputFromUser("date of event")
			t, err := dateParser(inputStr)
			if err != nil {
				displayError(ctx, "Wrong formatting !")
			} else {
				*date = time.Date(t.Year(), t.Month(), t.Day(), date.Hour(), date.Minute(), date.Second(), 0, time.Local)
				askAgain = false
			}
		}
	}
	askTime := func(date *time.Time) {
		askAgain := true
		for askAgain {
			inputStr := inputFromUser("begin time of event")
			t, err := timeParser(inputStr)
			if err != nil {
				displayError(ctx, "Wrong formatting !")
			} else {
				*date = time.Date(date.Year(), date.Month(), date.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.Local)
				askAgain = false
			}
		}
	}
	askEndTime := func(endDate *time.Time) {
		askAgain := true
		for askAgain {
			inputStr := inputFromUser("end time of event")
			t, err := timeParser(inputStr)
			if err != nil {
				displayError(ctx, "Wrong formatting !")
			} else {
				*endDate = time.Date(date.Year(), date.Month(), date.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.Local)
				if !endDate.After(date) {
					displayError(ctx, "End time cannot be before start time !")
				} else {
					askAgain = false
				}
			}
		}
	}
	askName := func(name *string) {
		*name = inputFromUser("name of event")
	}
	askCategory := func(category *string) {
		*category = inputFromUser("category of event")
	}

	if len(command) == 2 {
		// One argument given, we need to check which one is it : time, date or name
		date, err = timeParser(command[1])
		if err == nil { // we have our time
			isTimeSet = true // So we set this flag on
		} else {
			date, err = dateParser(command[1])
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

		errTime, errDate := buildDateFromTimeDate(command[1], command[2], &date)
		if errTime == nil && errDate == nil {
			isTimeSet = true
			isDateSet = true
		} else if errTime == nil && errDate != nil { // getting date failed
			isTimeSet = true
			category = command[2]
		} else {
			// date time
			// date category
			errTime, errDate := buildDateFromDateTime(command[1], command[2], &date)
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

		errTime, errDate := buildDateFromTimeDate(command[1], command[2], &date)
		if errTime == nil && errDate == nil {
			// time date category
			isTimeSet = true
			isDateSet = true

			endDate, err = timeParser(command[3])
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
			errTime, errDate := buildDateFromDateTime(command[1], command[2], &date)
			if errTime == nil && errDate == nil {
				// date time category
				isTimeSet = true
				isDateSet = true
				// Check if date time endDate
				endDate, err = timeParser(command[3])
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

	color := confGetColorFromName(category, ctx.configuration)
	displayOk(ctx, "Adding event "+name+" of category "+category+" starting "+date.Format("2006-01-02")+" at "+date.Format("15:04")+" until "+endDate.Format("15:04"))
	_, err = insertActivity(name, color, date, endDate, ctx.srv)
	if err != nil {
		displayError(ctx, err.Error())
	}
	return err
}

// Print usage
func helpCommand(command Command, ctx *gogendaContext) {
	prefix := ""
	if !ctx.isShell {
		prefix = " gogenda"
	}
	specificHelp := ""
	if len(command) == 2 {
		specificHelp = command[1]
	}
	if specificHelp == "" {
		displayInfoHeading(ctx, "== GoGenda ==")
		fmt.Println(" GoGenda helps you keep track of your activities")
		fmt.Println(" Type Gogenda help (command) to have more help for a specific command")
		displayInfoHeading(ctx, " = Commands = ")
		fmt.Println("")
		if !ctx.isShell {
			fmt.Println(" gogenda shell - Launch the shell UI")
		}
		for _, category := range ctx.configuration.Categories {
			fmt.Println(prefix + " start " + category.Name + " - Add an event in " + category.Color)
		}
		fmt.Println(prefix + " stop - Stop the current activity")
		fmt.Println(prefix + " rename - Rename the current activity")
		fmt.Println(prefix + " delete - Delete the current activity")
		fmt.Println(prefix + " plan - shows events of the day. You can call it alone or with a date param.")
		fmt.Println(prefix + " add - add an event to the planning. You can call it alone or with some params.")
		fmt.Println(prefix + " help - shows the help")
		fmt.Println(prefix + " version - shows the current version")
	} else if strings.ToUpper(specificHelp) == "ADD" {
		fmt.Println(prefix + " add - add an event to the planning. You can call it alone or with some params.")
		fmt.Println("  | the program will ask you the remaining parameters of the event")
		fmt.Println("  | (time) ")
		fmt.Println("  | (time) (date)")
		fmt.Println("  | (time) (date) (endTime)")
		fmt.Println("  | (time) (date) (endTime) (category)")
		fmt.Println("  | (time) (date) (endTime) (category) (name...)")
		fmt.Println("  | (time) (date) (category)")
		fmt.Println("  | (time) (date) (category) (name...)")
		fmt.Println("  | (time) (category) (name)")
		fmt.Println("  | (date) ")
		fmt.Println("  | (date) (time)")
		fmt.Println("  | (date) (time) (endTime)")
		fmt.Println("  | (date) (time) (endTime) (category)")
		fmt.Println("  | (date) (time) (endTime) (category) (name...)")
		fmt.Println("  | (date) (time) (category)")
		fmt.Println("  | (date) (time) (category) (name...)")
		fmt.Println("  - (date) (category) (name)")
	} else if strings.ToUpper(specificHelp) == "PLAN" {
		fmt.Println(prefix + " plan - shows events of the day. You can call it alone or with a date param.")
		fmt.Println("  | The program will get you today's planning if you dont specify a param")
		fmt.Println("  - (date)")
		fmt.Println("  - (date) (nb of days)")
	}

	if specificHelp != "" {
		fmt.Println(" Param guide : (time) can be, case unsensitive, 'now', 'HH', 'HH:MM', 'HH:MM:SS'")
		fmt.Println("             | (date) can be, case unsensitive, 'yesterday', 'today', 'tomorrow', 'YYYY-MM-DD', 'YYYY/MM/DD', 'MM-DD', 'MM/DD'")
		fmt.Println("             | (category) is one of the one you declared in your config.json file, case unsensitive")
	}
}
