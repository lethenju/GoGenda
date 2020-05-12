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
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

// Add an event now
func startCommand(command []string, ctx *gogendaContext) (err error) {
	var nameOfEvent string
	color := confGetColorForName(command[1], ctx.configuration)
	if len(command) == 2 && color != "blue" {
		fmt.Print(command)
		fmt.Print("Enter name of event :")
		scanner := bufio.NewScanner(os.Stdin)
		if !scanner.Scan() {
			return
		}
		nameOfEvent = scanner.Text()

		if ctx.activity.Id != "" {
			// Stop the current activity
			err = stopActivity(ctx.activity, ctx.srv)
			if err != nil {
				displayError(ctx, "There was an issue deleting the current event.")
			}
		}
	} else if len(command) == 2 {
		nameOfEvent = strings.Join(command[1:], " ")
	} else {
		nameOfEvent = strings.Join(command[2:], " ")
	}

	*ctx.activity, err = insertActivity(nameOfEvent, color, time.Now(), time.Now().Add(30*time.Minute), ctx.srv)

	if err != nil {
		return err
	}
	displayOk(ctx, "Successfully added activity ! ")
	return nil
}

func stopCommand(ctx *gogendaContext) (err error) {
	if ctx.activity.Id == "" {
		// Nothing to stop
		return errors.New("Nothing to stop")
	}

	duration, err := getDuration(ctx.activity)
	if err != nil {
		return err
	}
	displayInfo(ctx, "The activity '"+ctx.activity.Summary+"' lasted "+duration)
	if err != nil {
		return err
	}
	err = stopActivity(ctx.activity, ctx.srv)

	displayOk(ctx, "Successfully stopped the activity ! I hope it went well ")
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

func renameCommand(command []string, ctx *gogendaContext) (err error) {
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

func planCommand(command []string, ctx *gogendaContext) (err error) {
	// Get plan of all day
	begin := time.Now()
	begin = time.Date(begin.Year(), begin.Month(), begin.Day(), 0, 0, 0, 0, time.Local)
	if len(command) > 1 {
		begin, err = dateParser(command[1])
	}

	end := time.Date(begin.Year(), begin.Month(), begin.Day(), 23, 59, 59, 0, time.Local)

	cals, err := getActivitiesBetweenDates(begin.Format(time.RFC3339), end.Format(time.RFC3339), ctx.srv)
	displayInfoHeading(ctx, " Events of "+begin.Format(time.RFC822))
	if cals == nil {
		displayError(ctx, "Error")
		return err
	}
	events := cals.Items
	for _, event := range events {
		beginTime, _ := time.Parse(time.RFC3339, event.Start.DateTime)
		endTime, _ := time.Parse(time.RFC3339, event.End.DateTime)

		displayOk(ctx, " [ "+beginTime.Format(time.RFC822)+" -> "+endTime.Format(time.RFC822)+" ] : "+event.Summary)
	}
	return err
}

// Add an event sometime
// If you want to add it now, you better use startCommand
func addCommand(command []string, ctx *gogendaContext) (err error) {
	var date time.Time
	var endDate time.Time

	var name string
	var category string
	if len(command) == 1 {
		// No arguments given, we're gonna ask the user everything
		askAgain := true
		for askAgain {
			inputStr := inputFromUser("date of event")
			date, err = dateParser(inputStr)
			if err != nil {
				displayError(ctx, "Wrong formatting !")
			} else {
				askAgain = false
			}
		}

		askAgain = true
		for askAgain {
			inputStr := inputFromUser("begin time of event")
			t, err := timeParser(inputStr)
			if err != nil {
				displayError(ctx, "Wrong formatting !")
			} else {
				date = time.Date(date.Year(), date.Month(), date.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.Local)
				askAgain = false
			}
		}

		askAgain = true
		for askAgain {
			inputStr := inputFromUser("end time of event")
			t, err := timeParser(inputStr)
			if err != nil {
				displayError(ctx, "Wrong formatting !")
			} else {
				endDate = time.Date(date.Year(), date.Month(), date.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.Local)
				if !endDate.After(date) {
					displayError(ctx, "End time cannot be before start time !")
				} else {
					askAgain = false
				}
			}
		}
		name = inputFromUser("name of event")
		category = inputFromUser("category of event")
	}

	color := confGetColorForName(category, ctx.configuration)
	_, err = insertActivity(name, color, date, endDate, ctx.srv)
	if err != nil {
		displayError(ctx, err.Error())
	}
	return err
}

// Print usage
func helpCommand(ctx *gogendaContext) {
	displayInfoHeading(ctx, "== GoGenda ==")
	fmt.Println(" GoGenda helps you keep track of your activities")
	displayInfoHeading(ctx, " = Commands = ")
	prefix := ""
	if !ctx.isShell {
		prefix = " gogenda"
	}
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
	fmt.Println(prefix + " plan (today / tommorow / yyyy-mm-dd / mm-dd) - shows events of the day")
	fmt.Println(prefix + " add - add an event to the planning")
	fmt.Println(prefix + " help - shows the help")
	fmt.Println(prefix + " version - shows the current version")
}
