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
 @Version : 0.1.4
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
)

func startCommand(command []string, ctx *gogendaContext) (err error) {
	var nameOfEvent string
	if len(command) == 2 {
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
	} else {
		nameOfEvent = strings.Join(command[2:], " ")
	}
	ourCategory := ConfigCategory{Name: "default", Color: "blue"}
	for _, category := range ctx.configuration.Categories {
		if strings.ToUpper(command[1]) == category.Name {
			ourCategory = category
		}
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

func helpCommand(ctx *gogendaContext) {
	displayInfoHeading(ctx, "== GoGenda ==")
	fmt.Println(" GoGenda helps you keep track of your activities")
	displayInfoHeading(ctx, " = Commands = ")
	fmt.Println("")
	for _, category := range ctx.configuration.Categories {
		fmt.Println(" START " + category.Name + " - Add an event in " + category.Color)
	}
	fmt.Println(" STOP - Stop the current activity")
	fmt.Println(" RENAME - Rename the current activity")
	fmt.Println(" DELETE - Delete the current activity")
}
