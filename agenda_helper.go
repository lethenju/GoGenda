package main

import (
	"time"

	"google.golang.org/api/calendar/v3"
)

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
	activity.Id = ""
	return err
}

func deleteActivity(activity *calendar.Event, srv *calendar.Service) (err error) {
	call := srv.Events.Delete("primary", activity.Id)
	err = call.Do()
	activity.Id = ""
	return err
}
func renameActivity(activity *calendar.Event, text string, srv *calendar.Service) (err error) {
	var edtEnd calendar.EventDateTime
	edtEnd.DateTime = time.Now().Format(time.RFC3339)
	activity.Summary = text
	call := srv.Events.Update("primary", activity.Id, activity)
	_, err = call.Do()
	return err
}
