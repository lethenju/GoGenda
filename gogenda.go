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
 @Version : 0.1.1
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

	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

const version = "0.1.1"

func commandHandler(command []string, activity *calendar.Event, srv *calendar.Service) (err error) {

	switch strings.ToUpper(command[0]) {
	case "START":
		var nameOfEvent string
		if len(command) == 2 {
			fmt.Print("Enter name of event :  ")
			scanner := bufio.NewScanner(os.Stdin)
			if !scanner.Scan() {
				return
			}
			nameOfEvent = scanner.Text()

			if activity.Id != "" {
				// Stop the current activity
				err = stopActivity(activity, srv)
				if err != nil {
					fmt.Println("There was an issue deleting the current event.")
				}
			}
		} else {
			nameOfEvent = strings.Join(command[2:], " ")
		}

		switch strings.ToUpper(command[1]) {
		case "WORK":
			*activity, err = insertActivity(nameOfEvent, "red", srv)
			break
		case "ORGA":
			*activity, err = insertActivity(nameOfEvent, "yellow", srv)
			break
		case "LUNCH":
			*activity, err = insertActivity(nameOfEvent, "purple", srv)
			break
		default:
			return errors.New("I didnt recognised this activity")
		}
		if err != nil {
			return err
		}
		fmt.Println("Successfully added activity ! ")
		break
	case "STOP":
		if activity.Id == "" {
			// Nothing to stop
			return errors.New("Nothing to stop")
		}

		startTime, err := time.Parse(time.RFC3339, activity.Start.DateTime)
		if err != nil {
			return err
		}
		duration := time.Since(startTime)
		fmt.Println("The activity '" + activity.Summary + "' lasted " + duration.Truncate(time.Second).String())
		if err != nil {
			return err
		}
		err = stopActivity(activity, srv)

		fmt.Println("Successfully stopped the activity ! I hope it went well ")
		break
	case "RENAME":
		if activity.Id == "" {
			return errors.New("You dont have a current activity to rename")
		}
		fmt.Print("Enter name of event :  ")
		scanner := bufio.NewScanner(os.Stdin)
		if !scanner.Scan() {
			return
		}
		nameOfEvent := scanner.Text()

		err = renameActivity(activity, nameOfEvent, srv)
		if err != nil {
			return err
		}
		fmt.Println("Successfully renamed the activity")
		break
	case "DELETE":
		if activity.Id == "" {
			// Nothing to stop
			return errors.New("Nothing to delete")
		}
		err = deleteActivity(activity, srv)
		if err != nil {
			return err
		}
		fmt.Println("Successfully deleted the activity ! I hope it went well ")
		break
	case "HELP":
		fmt.Println("== GoGenda ==")
		fmt.Println(" GoGenda helps you keep track of your activities")
		fmt.Println(" = Commands = ")
		fmt.Println(" START WORK - Start a work related activity")
		fmt.Println(` START ORGA - Start a organisation related activity - 
		Reading articles, answering mails etc`)
		fmt.Println(" START LUNCH - Start a lunch related activity")
		fmt.Println(" STOP - Stop the current activity")
		fmt.Println(" RENAME - Rename the current activity")
		fmt.Println(" DELETE - Delete the current activity")
	default:
		fmt.Println(command[0] + ": command not found")
	}

	return nil
}

func main() {
	fmt.Println("Welcome to GoGenda!")
	fmt.Println("Version number : " + version)
	runningFlag := true
	var currentActivity calendar.Event

	b, err := ioutil.ReadFile("/etc/gogenda/credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := calendar.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	for runningFlag {

		scanner := bufio.NewScanner(os.Stdin)
		var command []string
		for len(command) == 0 {
			if currentActivity.Id != "" {
				fmt.Print("[" + currentActivity.Summary + "]")
			}
			fmt.Print("> ")
			if !scanner.Scan() {
				return
			}
			userInput := scanner.Text()
			command = strings.Fields(userInput)
		}
		if strings.ToUpper(command[0]) == "EXIT" {
			println("See you later !")
			if currentActivity.Id != "" {
				stopActivity(&currentActivity, srv)
			}
			runningFlag = false
			break
		}
		res := commandHandler(command, &currentActivity, srv)
		if res != nil {
			println("There was an error " + res.Error())
		}
	}

}
