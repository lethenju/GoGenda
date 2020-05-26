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
package utilities

import (
	"encoding/json"
	"os"
)

// EventStored is the data of an event of plan show that is stored for performing actions on previously got events
type EventStored struct {

	// Name
	Name string `json:"name"`
	// CalendarID
	CalendarID string `json:"CalendarID"`
}

// Plan is the type of the stored plan with ID of events to be modified
type Plan struct {
	// Events
	Events []EventStored `json:"events"`
}

// LoadPlan loads the stored plan
func LoadPlan() (plan Plan, err error) {
	f, err := os.Open("/tmp/plan")
	if err != nil {
		return plan, err
	}
	defer f.Close()
	plan = Plan{}
	err = json.NewDecoder(f).Decode(&plan)
	return plan, err
}

// StorePlan saves the stored plan
func StorePlan(plan *Plan) (err error) {
	os.Remove("/tmp/plan")
	f, err := os.Create("/tmp/plan")
	if err != nil {
		return err
	}
	defer f.Close()
	json.NewEncoder(f).Encode(plan)
	return err
}
