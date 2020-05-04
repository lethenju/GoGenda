package main

import (
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

func insertActivity(name string, color string, srv *calendar.Service) (err error) {
	var newEvent calendar.Event
	newEvent.Start.DateTime = time.Now().Format(time.RFC3339)
	newEvent.ColorId = "11" // ???
	newEvent.Summary = name
	srv.Events.Insert("primary", &newEvent)
	return nil
}

func stopActivity(activity *calendar.Event, srv *calendar.Service) (err error) {
	activity.End.DateTime = time.Now().Format(time.RFC3339)
	srv.Events.Update("primary", activity.Id, activity)
	activity = nil
	return nil
}

func commandHandler(command []string, activity *calendar.Event, srv *calendar.Service) (err error) {

	switch command[0] {
	case "START":

		var nameOfEvent string
		fmt.Print("Enter name of event :  ")
		fmt.Scan(&nameOfEvent)

		if activity != nil {
			// Stop the current activity
			stopActivity(activity, srv)
		}
		switch command[0] {
		case "WORK":
			insertActivity(nameOfEvent, "red", srv)
			break
		case "VEILLE":
			insertActivity(nameOfEvent, "yellow", srv)
			break
		case "REPAS":
			insertActivity(nameOfEvent, "purple", srv)
			break
		default:
			return errors.New("I didnt recognised this activity")
		}
		print("Successfully added event ! Work hard! ")
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
		print("Successfully stopped the event ! I hope it went well ")

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

		var userInput string
		fmt.Print("> ")
		fmt.Scan(&userInput)
		userInput = strings.ToUpper(userInput)
		command := strings.Split(userInput, " ")
		fmt.Println(len(command))

		res := commandHandler(command, currentActivity, srv)
		if res != nil {
			println("There was an error " + res.Error())
		}
	}
}
