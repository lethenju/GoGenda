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
	"os/user"
	"strings"

	gogenda "github.com/lethenju/gogenda/internal"
	cmdOptions "github.com/lethenju/gogenda/internal/cmd_options"
	"github.com/lethenju/gogenda/internal/configuration"
	"github.com/lethenju/gogenda/internal/current_activity"
	"github.com/lethenju/gogenda/internal/gogendalib"
	"github.com/lethenju/gogenda/pkg/colors"
	api "github.com/lethenju/gogenda/pkg/google_agenda_api"
)

// Version of the software
const version = "0.2.1"

// Main entry point
func main() {
	usr, _ := user.Current()
	userDir := usr.HomeDir

	args := cmdOptions.Init()
	// Setup colors printing
	colors.SetupColors()
	// Connect to API
	srv, err := api.Connect()

	if err != nil {
		colors.DisplayError(err.Error())
		return
	}
	config, err := cmdOptions.GetStringOption("config")
	if err != nil {
		// Load default configuration
		config = userDir + "/.gogenda/config.json"
	}
	// Load user defined config (absolute path)
	err = configuration.LoadConfiguration(config)
	if err != nil {
		// Conf doesnt exist
		colors.DisplayError("Could not open " + config)
	}
	if cmdOptions.IsOptionSet("help") {
		if len(args) > 0 {
			gogendalib.CommandHandler([]string{"HELP", args[0]}, srv, false)
		} else {
			gogendalib.CommandHandler([]string{"HELP"}, srv, false)
		}
		return
	}

	if len(args) > 0 {

		// For the other commands than start its obvious he/she is
		if strings.ToUpper(args[0]) != "START" {
			currentActivity, _ := api.GetLastEvent(srv)
			current_activity.SetCurrentActivity(&currentActivity)
		}
		err = gogendalib.CommandHandler(args, srv, false)
		if err != nil {
			colors.DisplayError("ERROR : " + err.Error())
		}
	} else if cmdOptions.GetNumberOfOptions() == 0 {
		// gogenda was called alone
		gogendalib.CommandHandler([]string{"HELP"}, srv, false)
	}
	// Launch shell based UI
	if cmdOptions.IsOptionSet("shell") {
		gogenda.Shell(srv, version)
		return
	}

}
