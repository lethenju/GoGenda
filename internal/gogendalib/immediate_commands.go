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
	"strings"
	"time"

	"github.com/lethenju/gogenda/internal/configuration"
	"github.com/lethenju/gogenda/internal/current_activity"
	"github.com/lethenju/gogenda/internal/utilities"
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
	current_activity.SetCurrentActivity(&currentActivity)

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

	current_activity.SetCurrentActivity(nil)

	colors.DisplayOk("Successfully stopped the activity ! I hope it went well ")
	return nil
}

func deleteCommand(srv *calendar.Service) (err error) {

	currentActivity, err := current_activity.GetCurrentActivity()
	if err != nil {
		return errors.New("Nothing to delete")
	}
	err = api.DeleteActivity(currentActivity, srv)
	if err != nil {
		return err
	}
	colors.DisplayOk("Successfully deleted the activity ! ")
	return nil
}

func renameCommand(command Command, srv *calendar.Service) (err error) {
	currentActivity, err := current_activity.GetCurrentActivity()
	if err != nil {
		return errors.New("Nothing to rename")
	}
	var nameOfEvent string
	if len(command) == 1 {
		nameOfEvent = utilities.InputFromUser("name of event")
	} else {
		nameOfEvent = strings.Join(command[1:], " ")
	}
	err = api.RenameActivity(currentActivity, nameOfEvent, srv)
	if err != nil {
		return err
	}
	colors.DisplayOk("Successfully renamed the activity")
	return nil
}
