package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func insertActivity(name string, color string, srv *calendar.Service) (activity calendar.Event, err error) {
	var newEvent calendar.Event
	var edtStart calendar.EventDateTime
	var edtEnd calendar.EventDateTime
	edtStart.DateTime = time.Now().Format(time.RFC3339)
	edtEnd.DateTime = time.Now().Add(30 * time.Minute).Format(time.RFC3339)
	newEvent.Start = &edtStart
	newEvent.End = &edtEnd
	switch color {
	case "red":
		newEvent.ColorId = "11"
		break
	case "yellow":
		newEvent.ColorId = "5"
		break
	case "purple":
		newEvent.ColorId = "3"
		break
	}
	newEvent.Summary = name
	call := srv.Events.Insert("primary", &newEvent)
	actualEvent, err := call.Do()
	newEvent.Id = actualEvent.Id
	return newEvent, err
}

func stopActivity(activity *calendar.Event, srv *calendar.Service) (err error) {
	var edtEnd calendar.EventDateTime
	edtEnd.DateTime = time.Now().Format(time.RFC3339)
	activity.End = &edtEnd
	call := srv.Events.Update("primary", activity.Id, activity)
	_, err = call.Do()
	activity = nil
	return err
}

func commandHandler(command []string, activity *calendar.Event, srv *calendar.Service) (err error) {

	switch command[0] {
	case "START":

		fmt.Print("Enter name of event :  ")
		scanner := bufio.NewScanner(os.Stdin)
		if !scanner.Scan() {
			return
		}
		nameOfEvent := scanner.Text()

		if activity != nil {
			// Stop the current activity
			err = stopActivity(activity, srv)
			if err != nil {
				fmt.Println("There was an issue deleting the current event.")
			}
		}
		switch command[1] {
		case "WORK":
			*activity, err = insertActivity(nameOfEvent, "red", srv)
			break
		case "VEILLE":
			*activity, err = insertActivity(nameOfEvent, "yellow", srv)
			break
		case "REPAS":
			*activity, err = insertActivity(nameOfEvent, "purple", srv)
			break
		default:
			return errors.New("I didnt recognised this activity")
		}
		if err != nil {
			return err
		}
		fmt.Println("Successfully added event ! Work hard! ")
		break
	case "STOP":
		if activity == nil {
			// Nothing to stop
			return errors.New("Nothing to stop")
		}
		err = stopActivity(activity, srv)
		if err != nil {
			return err
		}
		fmt.Println("Successfully stopped the event ! I hope it went well ")

		break
	}

	return nil
}

func main() {
	fmt.Println("Welcome to GoGenda!")
	runningFlag := true
	var currentActivity *calendar.Event

	b, err := ioutil.ReadFile("credentials.json")
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

	t := time.Now().Format(time.RFC3339)

	events, err := srv.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}
	if len(events.Items) > 0 {
		currentActivity = events.Items[0]
	}

	for runningFlag {

		scanner := bufio.NewScanner(os.Stdin)
		var command []string
		for len(command) == 0 {
			fmt.Print("> ")
			if !scanner.Scan() {
				return
			}
			userInput := scanner.Text()
			command = strings.Fields(userInput)
			userInput = strings.ToUpper(userInput)
		}

		res := commandHandler(command, currentActivity, srv)
		if res != nil {
			println("There was an error " + res.Error())
		}
	}

}
