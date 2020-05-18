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
 @Version : 0.2.0
 @Author : Julien LE THENO
 =============================================
*/
package main

import (
	"os"
	"os/user"
	"strings"

	"github.com/lethenju/gogenda/pkg/colors"
	"github.com/lethenju/gogenda/pkg/google_agenda_api"

	"google.golang.org/api/calendar/v3"
)

// Version of the software
const version = "0.2.0"

// GogendaContext The gogendaContext type centralises every needed data of the application.
type GogendaContext struct {
	// Current activity
	// The "service" of the calendar API, to able us to call API methods in Google Calendar's endpoint
	srv *calendar.Service
	// if we're on shell mode or not
	isShell bool
}

// The commandHandler takes the command in parameter and dispatchs it to the different command methods in command.go
func commandHandler(command []string, ctx *gogendaContext) (err error) {
	// Our command name is in the first argument
	switch strings.ToUpper(command[0]) {
	// Start an event
	case "START":
		err = startCommand(command, ctx)
		if err != nil {
			return err
		}
		break
	case "STOP":
		// Stop an event
		err = stopCommand(ctx)
		if err != nil {
			return err
		}
		break
	case "RENAME":
		// Renames an event
		err = renameCommand(command, ctx)
		if err != nil {
			return err
		}
		break
	case "DELETE":
		// Deletes an event
		err = deleteCommand(ctx)
		if err != nil {
			return err
		}
		break
	case "PLAN":
		// Show the plan of the date (or today if no date)
		err = planCommand(command, ctx)
		if err != nil {
			return err
		}
		break
	case "ADD":
		// add an event to the calendar at a specific date
		err = addCommand(command, ctx)
		if err != nil {
			return err
		}
		break
	case "HELP":
		// Show help
		helpCommand(command, ctx)
	case "VERSION":
		// Display version
		displayInfo(ctx, "Gogenda (MIT) Version : "+version)
	default:
		displayError(ctx, command[0]+": command not found")

	}

	return nil
}

// Main entry point
func main() {
	usr, _ := user.Current()
	userDir := usr.HomeDir

	var ctx gogendaContext
	// Setup colors printing
	colors.SetupColors()
	// Connect to API
	ctx.srv, err = google_agenda_api.Connect()
	if err != nil {
		displayError(err)
		return
	}

	ctx.activity = &currentActivity

	config, err := LoadConfiguration(userDir + "/.gogenda/config.json")
	if err != nil {
		// Conf doesnt exist
		displayError("Could not open ~/.gogenda/config.json")
	}
	args := os.Args
	if len(args) > 1 {
		// Launch shell based UI
		if strings.ToUpper(args[1]) == "SHELL" || strings.ToUpper(args[1]) == "-SH" {
			gogenda.shell(&ctx)
			return
		}
		ctx.isShell = false

		if strings.ToUpper(args[1]) == "HELP" || strings.ToUpper(args[1]) == "--HELP" {
			helpCommand([]string{}, &ctx)
			return
		}

		// For the other commands than start its obvious he/she is
		if strings.ToUpper(args[1]) != "START" {
			currentActivity, _ = getLastEvent(&ctx)
		}
		err := commandHandler(args[1:], &ctx)
		if err != nil {
			displayError(&ctx, "ERROR : "+err.Error())
		}
		return
	}
	helpCommand([]string{}, &ctx)
}
