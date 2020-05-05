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

func getDuration(activity *calendar.Event) (string, error) {

	startTime, err := time.Parse(time.RFC3339, activity.Start.DateTime)
	if err != nil {
		return "", err
	}
	duration := time.Since(startTime)
	return duration.Truncate(time.Second).String(), nil
}
