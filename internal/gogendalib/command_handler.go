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
	"strings"

	"github.com/lethenju/gogenda/pkg/colors"
	"google.golang.org/api/calendar/v3"
)

// CommandHandler takes the command in parameter and dispatchs it to the different command methods in command.go
func CommandHandler(command []string, srv *calendar.Service, isShell bool) (err error) {
	// Our command name is in the first argument
	switch strings.ToUpper(command[0]) {
	// Start an event
	case "START":
		err = startCommand(command, srv)
		if err != nil {
			return err
		}
		break
	case "STOP":
		// Stop an event
		err = stopCommand(srv)
		if err != nil {
			return err
		}
		break
	case "RENAME":
		// Renames an event
		err = renameCommand(command, srv)
		if err != nil {
			return err
		}
		break
	case "DELETE":
		// Deletes an event
		err = deleteCommand(srv)
		if err != nil {
			return err
		}
		break
	case "PLAN":
		// Show the plan of the date (or today if no date)
		err = planCommand(command, srv)
		if err != nil {
			return err
		}
		break
	case "ADD":
		// add an event to the calendar at a specific date
		err = addCommand(command, srv)
		if err != nil {
			return err
		}
		break
	case "STATS":
		// add an event to the calendar at a specific date
		err = statsCommand(command, srv)
		if err != nil {
			return err
		}
		break
	case "HELP":
		// Show help
		helpCommand(command, isShell)
	default:
		colors.DisplayError(command[0] + ": command not found")

	}

	return nil
}
