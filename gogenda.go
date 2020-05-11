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
 @Version : 0.1.6
 @Author : Julien LE THENO
 =============================================
*/
package main

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"runtime"
	"strings"
	"time"

	"google.golang.org/api/calendar/v3"
)

const version = "0.1.6"

type gogendaContext struct {
	activity      *calendar.Event
	srv           *calendar.Service
	colors        colors
	configuration Config
}

func commandHandler(command []string, ctx *gogendaContext) (err error) {

	switch strings.ToUpper(command[0]) {
	case "START":
		err = startCommand(command, ctx)
		if err != nil {
			return err
		}
		break
	case "STOP":
		err = stopCommand(ctx)
		if err != nil {
			return err
		}
		break
	case "RENAME":
		err = renameCommand(command, ctx)
		if err != nil {
			return err
		}
		break
	case "DELETE":
		err = deleteCommand(ctx)
		if err != nil {
			return err
		}
		break
	case "PLAN":
		err = planCommand(command, ctx)
		if err != nil {
			return err
		}
		break
	case "HELP":
		helpCommand(ctx)
	case "VERSION":
		displayInfo(ctx, "Gogenda (MIT) Version : "+version)
	default:
		displayError(ctx, command[0]+": command not found")

	}

	return nil
}

func getLastEvent(ctx *gogendaContext) (calendar.Event, error) {

	var selectedEvent calendar.Event

	t := time.Now().Format(time.RFC3339)
	events, err := ctx.srv.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(time.Now().Add(-12 * time.Hour).Format(time.RFC3339)).TimeMax(t).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		displayError(ctx, "ERROR : "+err.Error())
		return selectedEvent, err
	}

	var oldTime calendar.EventDateTime
	oldTime.DateTime = time.RFC3339
	selectedEvent.Start = &oldTime
	for _, event := range events.Items {
		if event.Start.DateTime > selectedEvent.Start.DateTime {
			selectedEvent = *event
		}
	}
	return selectedEvent, nil
}

func shell(ctx *gogendaContext) {

	runningFlag := true
	displayInfoHeading(ctx, "Welcome to GoGenda!")
	displayInfo(ctx, "Version number : "+version)
	lastEvent, err := getLastEvent(ctx)
	if err == nil && lastEvent.Id != "" {
		fmt.Println("Last event : " + lastEvent.Summary)
		fmt.Println("Are you still doing that ? (y/n)")
		userInput := ""
		for userInput != "y" && userInput != "n" {
			fmt.Scan(&userInput)
		}
		if userInput == "y" {
			ctx.activity = &lastEvent
		}
	}
	var userInput string

	scanner := bufio.NewScanner(os.Stdin)

	if runtime.GOOS == "windows" {
		// Scan twice on windows because scanner is not empty at startup
		if !scanner.Scan() {
			return
		}
		userInput = scanner.Text()
	}
	for runningFlag {

		var command []string
		for len(command) == 0 {
			if ctx.activity.Id != "" {
				fmt.Print("[ ")
				displayOkNoNL(ctx, ctx.activity.Summary+" ")
				duration, err := getDuration(ctx.activity)
				if err != nil {
					displayError(ctx, "ERROR : "+err.Error())
				}
				displayInfoNoNL(ctx, duration)

				fmt.Print(" ]")
			}
			fmt.Print("> ")
			if !scanner.Scan() {
				return
			}

			userInput = scanner.Text()
			command = strings.Fields(userInput)
		}
		if strings.ToUpper(command[0]) == "EXIT" {
			fmt.Println("See you later !")
			if ctx.activity.Id != "" {
				stopActivity(ctx.activity, ctx.srv)
			}
			runningFlag = false
			break
		}
		res := commandHandler(command, ctx)
		if res != nil {
			displayError(ctx, "ERROR : "+res.Error())
		}
	}

}

func main() {

	usr, _ := user.Current()
	userDir := usr.HomeDir

	var ctx gogendaContext
	setupColors(&ctx)

	connect(&ctx)
	var currentActivity calendar.Event

	ctx.activity = &currentActivity

	err := LoadConfiguration(userDir+"/.gogenda/config.json", &ctx)
	if err != nil {
		// Conf doesnt exist
		displayError(&ctx, "Could not open ~/.gogenda/config.json")
	}
	args := os.Args
	if len(args) > 1 {
		// Launch shell based UI
		if strings.ToUpper(args[1]) == "SHELL" || strings.ToUpper(args[1]) == "-SH" {
			shell(&ctx)
			return
		}

		if strings.ToUpper(args[1]) == "HELP" || strings.ToUpper(args[1]) == "--HELP" {
			usageCommand(&ctx)
			return
		}

		// For the other commands than start its obvious he/she is
		if strings.ToUpper(args[1]) != "START" {
			currentActivity, _ = getLastEvent(&ctx)
		}
		// for the START command :
		// Need to make sure if the user thinks he/she's still going under the last event or not
		// For now, lets suppose he/she doesnt. He/she should call stop first then start

		err := commandHandler(args[1:], &ctx)
		if err != nil {
			displayError(&ctx, "ERROR : "+err.Error())
		}
		return
	}
	usageCommand(&ctx)
}
