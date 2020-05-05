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
 @Version : 0.1.3
 @Author : Julien LE THENO
 =============================================
*/
package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

const version = "0.1.3"

type colors struct {
	colorInfo        *color.Color
	colorInfoHeading *color.Color
	colorOk          *color.Color
	colorError       *color.Color
}

type gogendaContext struct {
	activity *calendar.Event
	srv      *calendar.Service
	colors   colors
}

func displayInfo(ctx *gogendaContext, str string) {
	ctx.colors.colorInfo.Print(str)
	fmt.Println()
}
func displayInfoNoNL(ctx *gogendaContext, str string) {
	ctx.colors.colorInfo.Print(str)
}
func displayInfoHeading(ctx *gogendaContext, str string) {
	ctx.colors.colorInfoHeading.Print(str)
	fmt.Println() // Because the newline keeps the background color, we need to do it separately
}
func displayError(ctx *gogendaContext, str string) {
	ctx.colors.colorError.Print(str)
	fmt.Println()
}
func displayOk(ctx *gogendaContext, str string) {
	ctx.colors.colorOk.Println(str)
}
func displayOkNoNL(ctx *gogendaContext, str string) {
	ctx.colors.colorOk.Print(str)
}

func getDuration(activity *calendar.Event) (string, error) {

	startTime, err := time.Parse(time.RFC3339, activity.Start.DateTime)
	if err != nil {
		return "", err
	}
	duration := time.Since(startTime)
	return duration.Truncate(time.Second).String(), nil
}

func commandHandler(command []string, ctx *gogendaContext) (err error) {

	switch strings.ToUpper(command[0]) {
	case "START":
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

		switch strings.ToUpper(command[1]) {
		case "WORK":
			*ctx.activity, err = insertActivity(nameOfEvent, "red", ctx.srv)
			break
		case "ORGA":
			*ctx.activity, err = insertActivity(nameOfEvent, "yellow", ctx.srv)
			break
		case "LUNCH":
			*ctx.activity, err = insertActivity(nameOfEvent, "purple", ctx.srv)
			break
		default:
			return errors.New("I didnt recognised this activity")
		}
		if err != nil {
			return err
		}
		displayOk(ctx, "Successfully added activity ! ")
		break
	case "STOP":
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
		break
	case "RENAME":
		if ctx.activity.Id == "" {
			return errors.New("You dont have a current activity to rename")
		}
		fmt.Print(ctx, "Enter name of event :  ")
		scanner := bufio.NewScanner(os.Stdin)
		if !scanner.Scan() {
			return
		}
		nameOfEvent := scanner.Text()

		err = renameActivity(ctx.activity, nameOfEvent, ctx.srv)
		if err != nil {
			return err
		}
		displayOk(ctx, "Successfully renamed the activity")
		break
	case "DELETE":
		if ctx.activity.Id == "" {
			// Nothing to stop
			return errors.New("Nothing to delete")
		}
		err = deleteActivity(ctx.activity, ctx.srv)
		if err != nil {
			return err
		}
		displayOk(ctx, "Successfully deleted the activity ! ")
		break
	case "HELP":
		displayInfoHeading(ctx, "== GoGenda ==")
		fmt.Println(" GoGenda helps you keep track of your activities")
		displayInfoHeading(ctx, " = Commands = ")
		fmt.Println("")
		fmt.Println(" START WORK - Start a work related activity")
		fmt.Println(` START ORGA - Start a organisation related activity - 
		Reading articles, answering mails etc`)
		fmt.Println(" START LUNCH - Start a lunch related activity")
		fmt.Println(" STOP - Stop the current activity")
		fmt.Println(" RENAME - Rename the current activity")
		fmt.Println(" DELETE - Delete the current activity")
	default:
		displayError(ctx, command[0]+": command not found")

	}

	return nil
}

func main() {

	var ctx gogendaContext

	ctx.colors.colorInfo = color.New(color.FgBlue).Add(color.BgWhite)
	ctx.colors.colorInfoHeading = color.New(color.FgWhite).Add(color.BgBlue)
	ctx.colors.colorOk = color.New(color.FgGreen)
	ctx.colors.colorError = color.New(color.FgWhite).Add(color.BgHiRed)

	displayInfoHeading(&ctx, "Welcome to GoGenda!")
	displayInfo(&ctx, "Version number : "+version)
	runningFlag := true
	var currentActivity calendar.Event
	ctx.activity = &currentActivity
	b, err := ioutil.ReadFile("/etc/gogenda/credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	ctx.srv, err = calendar.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	for runningFlag {

		scanner := bufio.NewScanner(os.Stdin)
		var command []string
		for len(command) == 0 {
			if currentActivity.Id != "" {
				fmt.Print("[ ")
				displayOkNoNL(&ctx, currentActivity.Summary+" ")
				duration, err := getDuration(ctx.activity)
				if err != nil {
					displayError(&ctx, "ERROR : "+err.Error())
				}
				displayInfoNoNL(&ctx, duration)

				fmt.Print(" ]")
			}
			fmt.Print("> ")
			if !scanner.Scan() {
				return
			}
			userInput := scanner.Text()
			command = strings.Fields(userInput)
		}
		if strings.ToUpper(command[0]) == "EXIT" {
			fmt.Println("See you later !")
			if currentActivity.Id != "" {
				stopActivity(&currentActivity, ctx.srv)
			}
			runningFlag = false
			break
		}
		res := commandHandler(command, &ctx)
		if res != nil {
			displayError(&ctx, "ERROR : "+res.Error())
		}
	}

}
