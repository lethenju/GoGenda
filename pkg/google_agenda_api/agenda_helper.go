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
package google_agenda_api

import (
	"errors"
	"time"

	"google.golang.org/api/calendar/v3"
)

// InsertActivity : Inserts an activity in the agenda
// with the name of the event and the color of the event you want, the start and end time
// colors can be : "red", "yellow", "purple", "orange", "blue"
// Also give a pointer the the calendar service in order to send the api.
// It will return, if it succeeds, the event created, and an error code in case it fails.
func InsertActivity(name string, color string, beginTime time.Time, endTime time.Time, srv *calendar.Service) (activity calendar.Event, err error) {
	var newEvent calendar.Event
	var edtStart calendar.EventDateTime
	var edtEnd calendar.EventDateTime
	edtStart.DateTime = beginTime.Format(time.RFC3339)
	edtEnd.DateTime = endTime.Format(time.RFC3339)
	newEvent.Start = &edtStart
	newEvent.End = &edtEnd
	// 1 is lavender
	// 2 is green (sauge)
	// 3 is purple
	// 4 is rose
	// 5 is yellow
	// 6 is orange
	// 7 is blue
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
	case "orange":
		newEvent.ColorId = "6"
		break
	case "blue":
		newEvent.ColorId = "7"
		break
	}
	// No necessary default case as ColorId doesnt have to be set
	newEvent.Summary = name
	call := srv.Events.Insert("primary", &newEvent)
	actualEvent, err := call.Do()
	newEvent.Id = actualEvent.Id
	return newEvent, err
}

// StopActivity : Stops the current activity : actually update the end time of the activity in parameters
// to be current time.
// Also give a pointer the the calendar service in order to send the api.
func StopActivity(activity *calendar.Event, srv *calendar.Service) (err error) {
	var edtEnd calendar.EventDateTime
	edtEnd.DateTime = time.Now().Format(time.RFC3339)
	activity.End = &edtEnd
	call := srv.Events.Update("primary", activity.Id, activity)
	_, err = call.Do()
	activity.Id = ""
	return err
}

// DeleteActivity : Deletes the activity given in parameters
// Also give a pointer the the calendar service in order to send the api.
func DeleteActivity(activity *calendar.Event, srv *calendar.Service) (err error) {
	call := srv.Events.Delete("primary", activity.Id)
	err = call.Do()
	activity.Id = ""
	return err
}

// DeleteActivityFromID : Deletes the activity related to the idgiven in parameters
// Also give a pointer the the calendar service in order to send the api.
func DeleteActivityFromID(EventID string, srv *calendar.Service) (err error) {
	call := srv.Events.Delete("primary", EventID)
	err = call.Do()
	return err
}

// MoveActivityFromID : Moves the activity with the datetime given in parareters
// Set the start time to the one in param, and stop time will be changed accordingly
// to keep the same duration
func MoveActivityFromID(EventID string, startTime time.Time, srv *calendar.Service) (err error) {
	request := srv.Events.Get("primary", EventID)
	event, err := request.Do()

	// Getting the duration of the activity
	oldEndTime, _ := time.Parse(time.RFC3339, event.End.DateTime)
	oldStartTime, _ := time.Parse(time.RFC3339, event.Start.DateTime)
	duration := oldEndTime.Sub(oldStartTime)

	event.Start.DateTime = startTime.Format(time.RFC3339)
	event.End.DateTime = startTime.Add(duration).Format(time.RFC3339)

	call := srv.Events.Update("primary", EventID, event)
	event, err = call.Do()
	// Todo check if it becomes the current event or not ?
	return err
}

// RenameActivity : Renames the activity given in parameters with the text parameter
// Also give a pointer the the calendar service in order to send the api.
func RenameActivity(activity *calendar.Event, text string, srv *calendar.Service) (err error) {
	activity.Summary = text
	call := srv.Events.Update("primary", activity.Id, activity)
	_, err = call.Do()
	return err
}

// GetActivitiesBetweenDates Retrieve a Events* list of events which occurs between the dates given in parameters (in format RFC3339)
// Also give a pointer the the calendar service in order to send the api.
func GetActivitiesBetweenDates(beginDate string, endDate string, srv *calendar.Service) (cals *calendar.Events, err error) {

	events, err := srv.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(beginDate).TimeMax(endDate).MaxResults(512).OrderBy("startTime").Do()
	return events, err
}

// GetDuration Retrieve the duration (now - startTime) of current event
func GetDuration(activity *calendar.Event) (string, error) {

	startTime, err := time.Parse(time.RFC3339, activity.Start.DateTime)
	if err != nil {
		return "", err
	}
	duration := time.Since(startTime)
	return duration.Truncate(time.Second).String(), nil
}

// GetColorIDFromColorName  TODO
func GetColorIDFromColorName(colorName string) (colorID string, err error) {
	switch colorName {
	case "red":
		return "11", nil
	case "yellow":
		return "5", nil
	case "purple":
		return "3", nil
	case "orange":
		return "6", nil
	case "blue":
		return "7", nil
	}
	return "", errors.New("Didnt find color")
}

// GetColorNameFromColorID TODO
func GetColorNameFromColorID(colorID string) (colorName string, err error) {
	switch colorID {
	case "11":
		return "red", nil
	case "5":
		return "yellow", nil
	case "3":
		return "purple", nil
	case "6":
		return "orange", nil
	case "7":
		return "blue", nil
	}
	return "7", errors.New("Didnt find color")
}

// GetLastEvent function gets the last event we set on google agenda today, in
// order to ask the user if he's still doing that task or not
// TODO : Use the newer getActivitiesBetweenDates instead
func GetLastEvent(srv *calendar.Service) (calendar.Event, error) {

	var selectedEvent calendar.Event

	t := time.Now().Format(time.RFC3339)
	events, err := srv.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(time.Now().Add(-12 * time.Hour).Format(time.RFC3339)).TimeMax(t).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		//displayError(ctx, "ERROR : "+err.Error())
		return selectedEvent, err
	}

	var oldTime calendar.EventDateTime
	oldTime.DateTime = time.RFC3339
	selectedEvent.Start = &oldTime
	for _, event := range events.Items {
		if event.Start.DateTime > selectedEvent.Start.DateTime {
			selectedEvent = *event
		}
	}
	return selectedEvent, nil
}

//GetStartDateForEventID returns the date of a event given its ID
func GetStartDateForEventID(ID string, srv *calendar.Service) (time.Time, error) {
	request := srv.Events.Get("primary", ID)
	event, err := request.Do()
	date, _ := time.Parse(time.RFC3339, event.Start.DateTime)
	return date, err
}

//GetEndDateForEventID returns the end date of a event given its ID
func GetEndDateForEventID(ID string, srv *calendar.Service) (time.Time, error) {
	request := srv.Events.Get("primary", ID)
	event, err := request.Do()
	date, _ := time.Parse(time.RFC3339, event.End.DateTime)
	return date, err
}
