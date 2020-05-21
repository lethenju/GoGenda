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
	"errors"
	"sort"
	"strconv"
	"time"

	cmdOptions "github.com/lethenju/gogenda/internal/cmd_options"
	"github.com/lethenju/gogenda/internal/configuration"
	"github.com/lethenju/gogenda/internal/utilities"
	"github.com/lethenju/gogenda/pkg/colors"
	api "github.com/lethenju/gogenda/pkg/google_agenda_api"

	"google.golang.org/api/calendar/v3"
)

func statsCommand(command Command, srv *calendar.Service) (err error) {
	// Get plan of all day
	begin := time.Now()
	begin = time.Date(begin.Year(), begin.Month(), begin.Day(), 0, 0, 0, 0, time.Local)
	if len(command) > 1 {
		begin, err = utilities.DateParser(command[1])
		if err != nil {
			return err
		}

	}

	nbDays := 1
	if len(command) > 2 {
		// Number of days to do
		nbDays, err = strconv.Atoi(command[2])
		if err != nil {
			return errors.New("Wrong argument '" + command[2] + "', should be a number")
		}
	}
	end := begin.Add(time.Duration(24*nbDays) * time.Hour)

	events, err := api.GetActivitiesBetweenDates(
		begin.Format(time.RFC3339),
		end.Format(time.RFC3339), srv)
	if err != nil {
		return err
	}
	items := events.Items

	if err != nil {
		return err
	}
	// sort by category
	sort.Slice(items, func(p, q int) bool {
		return items[p].ColorId < items[q].ColorId
	})

	lastColorCode := ""
	var total time.Duration
	for _, item := range items {
		if lastColorCode != item.ColorId {
			if lastColorCode != "" {
				colors.DisplayOk("      Total : " + total.String())
			}
			lastColorCode = item.ColorId
			total = 0
			// retrieve category
			colorName, _ := api.GetColorNameFromColorID(item.ColorId)
			category := configuration.GetNameFromColor(colorName)
			colors.DisplayInfoHeading("=== " + category + " ===")

		}
		startTime, _ := time.Parse(time.RFC3339, item.Start.DateTime)
		endTime, _ := time.Parse(time.RFC3339, item.End.DateTime)
		duration := endTime.Sub(startTime)
		total += duration
		if !cmdOptions.IsOptionSet("compact") {
			colors.DisplayOk(" [ " + startTime.Format("15:04") + " -> " + endTime.Format("15:04") + " ] " + duration.String() + " : " + item.Summary)
		}
	}
	colors.DisplayOk("      Total : " + total.String())

	return nil
}
