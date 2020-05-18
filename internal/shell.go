package gogenda

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/lethenju/gogenda/pkg/colors"
	"github.com/lethenju/gogenda/pkg/google_agenda_api"
)

// Gogenda can be called as a shell, to have a shell like environement for long periods of usage
func shell(ctx *GogendaContext) {
	ctx.isShell = true
	runningFlag := true

	colors.DisplayInfoHeading("Welcome to GoGenda!")
	colors.DisplayInfo("Version number : " + version)

	// Asking the user if he's still doing the last event on google agenda
	lastEvent, err := google_agenda_api.GetLastEvent(ctx.srv)
	if err == nil && lastEvent.Id != "" {
		fmt.Println("Last event : " + lastEvent.Summary)
		fmt.Println("Are you still doing that ? (y/n)")
		userInput := ""
		for userInput != "y" && userInput != "n" {
			fmt.Scan(&userInput)
		}
		if userInput == "y" {
			SetCurrentActivity(&lastEvent)
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

	// Main loop
	for runningFlag {

		var command []string
		for len(command) == 0 {
			act, err := GetCurrentActivity()
			if err != nil {
				fmt.Print("[ ")
				displayOkNoNL(ctx.activity.Summary + " ")
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
