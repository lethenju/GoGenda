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

func startCommand(command []string, ctx *gogendaContext) (err error) {
	var nameOfEvent string
	ourCategory := ConfigCategory{Name: "default", Color: "blue"}
	for _, category := range ctx.configuration.Categories {
		if strings.ToUpper(command[1]) == category.Name {
			ourCategory = category
		}
	}
	if len(command) == 2 && ourCategory.Name != "default" {
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

	*ctx.activity, err = insertActivity(nameOfEvent, ourCategory.Color, ctx.srv)

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
	end := time.Now()
	end = time.Date(begin.Year(), begin.Month(), begin.Day(), 23, 59, 59, 0, time.Local)
	if len(command) > 1 {
		if strings.ToUpper(command[1]) == "tommorow" {
			begin = begin.Add(24 * time.Hour)
			end = begin.Add(24 * time.Hour)
		}
	}

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

func helpCommand(ctx *gogendaContext) {
	displayInfoHeading(ctx, "== GoGenda ==")
	fmt.Println(" GoGenda helps you keep track of your activities")
	displayInfoHeading(ctx, " = Commands = ")
	fmt.Println("")
	for _, category := range ctx.configuration.Categories {
		fmt.Println(" start " + category.Name + " - Add an event in " + category.Color)
	}
	fmt.Println(" stop - Stop the current activity")
	fmt.Println(" rename - Rename the current activity")
	fmt.Println(" delete - Delete the current activity")
	fmt.Println(" help - shows the help")
	fmt.Println(" version - shows the current version")
}

func usageCommand(ctx *gogendaContext) {
	displayInfoHeading(ctx, "== GoGenda Usage ==")
	fmt.Println(" GoGenda helps you keep track of your activities")
	displayInfoHeading(ctx, " = Commands = ")
	fmt.Println("")
	fmt.Println(" gogenda shell - Launch the shell UI")
	for _, category := range ctx.configuration.Categories {
		fmt.Println(" gogenda start " + category.Name + " - Add an event in " + category.Color)
	}
	fmt.Println(" gogenda stop - Stop the current activity")
	fmt.Println(" gogenda rename - Rename the current activity")
	fmt.Println(" gogenda delete - Delete the current activity")
	fmt.Println(" gogenda help - shows the help")
	fmt.Println(" gogenda version - shows the current version")
}
