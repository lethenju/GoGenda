package gogenda

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"

	"google.golang.org/api/calendar/v3"

	"github.com/lethenju/gogenda/internal/current_activity"
	"github.com/lethenju/gogenda/internal/gogendalib"
	"github.com/lethenju/gogenda/pkg/colors"
	api "github.com/lethenju/gogenda/pkg/google_agenda_api"
)

//Shell : Gogenda can be called as a shell, to have a shell like environement for long periods of usage
func Shell(srv *calendar.Service, version string) {
	runningFlag := true

	colors.DisplayInfoHeading("Welcome to GoGenda!")
	colors.DisplayInfo("Version number : " + version)

	// Asking the user if he's still doing the last event on google agenda
	lastEvent, err := api.GetLastEvent(srv)
	if err == nil && lastEvent.Id != "" {
		fmt.Println("Last event : " + lastEvent.Summary)
		fmt.Println("Are you still doing that ? (y/n)")
		userInput := ""
		for userInput != "y" && userInput != "n" {
			fmt.Scan(&userInput)
		}
		if userInput == "y" {
			current_activity.SetCurrentActivity(&lastEvent)
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
			act, err := current_activity.GetCurrentActivity()
			if err != nil {
				fmt.Print("[ ")
				colors.DisplayOkNoNL(act.Summary + " ")
				duration, err := api.GetDuration(act)
				if err != nil {
					colors.DisplayError("ERROR : " + err.Error())
				}
				colors.DisplayInfoNoNL(duration)

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
			currentActivity, err := current_activity.GetCurrentActivity()
			if err == nil {
				api.StopActivity(currentActivity, srv)
			}
			runningFlag = false
			break
		}
		res := gogendalib.CommandHandler(command, srv, true)
		if res != nil {
			colors.DisplayError("ERROR : " + res.Error())
		}
	}

}
