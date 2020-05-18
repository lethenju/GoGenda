package gogenda

import (
	"errors"

	"google.golang.org/api/calendar/v3"
)

var currentActivity *calendar.Event

// GetCurrentActivity Returns the currentActivity
func GetCurrentActivity() (*calendar.Event, error) {
	if currentActivity == nil {
		return nil, errors.New("Current Activity is not defined")
	}
	return currentActivity, nil
}

// SetCurrentActivity sets the current activity
func SetCurrentActivity(activity *calendar.Event) {
	currentActivity = activity
}
