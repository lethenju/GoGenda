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
	"fmt"
	"os"
	"os/user"
	"strings"

	"google.golang.org/api/calendar/v3"
)

const version = "0.1.4"

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
	case "HELP":
		helpCommand(ctx)
	default:
		displayError(ctx, command[0]+": command not found")

	}

	return nil
}

func main() {
	usr, _ := user.Current()
	userDir := usr.HomeDir

	var ctx gogendaContext
	setupColors(&ctx)

	displayInfoHeading(&ctx, "Welcome to GoGenda!")
	displayInfo(&ctx, "Version number : "+version)
	runningFlag := true
	var currentActivity calendar.Event
	ctx.activity = &currentActivity

	LoadConfiguration(userDir+"/.gogenda/config.json", &ctx)

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
