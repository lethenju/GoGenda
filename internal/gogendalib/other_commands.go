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
	"fmt"
	"strings"

	"github.com/lethenju/gogenda/internal/configuration"
	"github.com/lethenju/gogenda/pkg/colors"
)

// Print usage
func helpCommand(command Command, isShell bool) {
	prefix := ""
	if !isShell {
		prefix = " gogenda"
	}
	specificHelp := ""
	if len(command) == 2 {
		specificHelp = command[1]
	}
	if specificHelp == "" {
		colors.DisplayInfoHeading("== GoGenda ==")
		fmt.Println(" GoGenda helps you keep track of your activities")
		fmt.Println(" Type Gogenda -h (command) to have more help for a specific command")
		fmt.Println("")
		if !isShell {
			colors.DisplayInfoHeading(" = Options = ")
			colors.DisplayOk("Important : options have to be used before command arguments !")
			fmt.Println(" gogenda -i              - Launch the shell UI")
			fmt.Println(" gogenda -h              - shows the help")
			fmt.Println(" gogenda -compact        - Have minimalist output")
			fmt.Println(" gogenda -config='path'  - Use a custom config file (absolute path only)")
			fmt.Println("")
		}
		colors.DisplayInfoHeading(" = Commands = ")

		config, _ := configuration.GetConfig()
		for _, category := range config.Categories {
			fmt.Println(prefix + " start " + category.Name + " - Add an event in " + category.Color)
		}
		fmt.Println(prefix + " stop - Stop the current activity")
		fmt.Println(prefix + " rename - Rename the current activity")
		fmt.Println(prefix + " delete - Delete the current activity")
		fmt.Println(prefix + " plan - shows events of one or several days. You can call it alone or with a date param.")
		fmt.Println(prefix + " stats - shows statistics about your time spent in each category")
		fmt.Println(prefix + " add - add an event to the planning. You can call it alone or with some params.")
		fmt.Println(prefix + " help - show gogenda help (add a command name if you want specific command help)")
	} else if strings.ToUpper(specificHelp) == "ADD" {
		fmt.Println(prefix + " add - add an event to the planning. You can call it alone or with some params.")
		fmt.Println("  | the program will ask you the remaining parameters of the event")
		fmt.Println("  | (time) ")
		fmt.Println("  | (time) (date)")
		fmt.Println("  | (time) (date) (endTime)")
		fmt.Println("  | (time) (date) (endTime) (category)")
		fmt.Println("  | (time) (date) (endTime) (category) (name...)")
		fmt.Println("  | (time) (date) (category)")
		fmt.Println("  | (time) (date) (category) (name...)")
		fmt.Println("  | (time) (category) (name)")
		fmt.Println("  | (date) ")
		fmt.Println("  | (date) (time)")
		fmt.Println("  | (date) (time) (endTime)")
		fmt.Println("  | (date) (time) (endTime) (category)")
		fmt.Println("  | (date) (time) (endTime) (category) (name...)")
		fmt.Println("  | (date) (time) (category)")
		fmt.Println("  | (date) (time) (category) (name...)")
		fmt.Println("  - (date) (category) (name)")
	} else if strings.ToUpper(specificHelp) == "PLAN" {
		fmt.Println(prefix + " plan - shows events of a one of several days. You can call it alone or with a date param.")
		fmt.Println("  | The program will get you today's planning if you don't specify a param")
		fmt.Println("  - (date)")
		fmt.Println("  - (date) (nb of days)")
	} else if strings.ToUpper(specificHelp) == "STATS" {
		fmt.Println(prefix + " stats - shows statistics about your time spent in each category")
		fmt.Println("  | The program will get you today's statistics if you don't specify a param")
		fmt.Println("  - (date)")
	}

	if specificHelp != "" {
		fmt.Println(" Param guide : (time) can be, case unsensitive, 'now', 'HH', 'HH:MM', 'HH:MM:SS'")
		fmt.Println("             | (date) can be, case unsensitive, 'yesterday', 'today', 'tomorrow', 'YYYY-MM-DD', 'YYYY/MM/DD', 'MM-DD', 'MM/DD'")
		fmt.Println("             | (category) is one of the one you declared in your config.json file, case unsensitive")
	}
}
