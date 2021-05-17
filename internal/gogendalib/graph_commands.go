package gogendalib

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	cmdOptions "github.com/lethenju/gogenda/internal/cmd_options"
	"github.com/lethenju/gogenda/internal/configuration"
	"github.com/lethenju/gogenda/internal/utilities"
	"github.com/lethenju/gogenda/pkg/colors"
	api "github.com/lethenju/gogenda/pkg/google_agenda_api"
	chart "github.com/wcharczuk/go-chart" //exposes "chart"
	"google.golang.org/api/calendar/v3"
)

func GraphCommand(command Command, srv *calendar.Service) (err error) {
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

	// We have to put a default 0 value that is not the zero value of item.ColorId (which is ""
	var total time.Duration
	var datesEvent []time.Time
	var durationWork []float64
	startDay1, _ := time.Parse(time.RFC3339, items[0].Start.DateTime)
	actualDay := startDay1.Day()
	for _, item := range items {
		startTime, _ := time.Parse(time.RFC3339, item.Start.DateTime)
		endTime, _ := time.Parse(time.RFC3339, item.End.DateTime)
		duration := endTime.Sub(startTime)
		if actualDay == startTime.Day() {

			// retrieve category
			colorName, _ := api.GetColorNameFromColorID(item.ColorId)
			category := configuration.GetNameFromColor(colorName)

			if category == "WORK" {
				total += duration
			}

		} else {
			colors.DisplayOk("CHANGEMENT DE JOUR")
			colors.DisplayOk("      Total : " + total.String())

			// changement de jour
			actualDay = startTime.Day()
			// Ajout du temps de travail de la journée dans le graph
			datesEvent = append(datesEvent, startTime)
			durationWork = append(durationWork, total.Minutes())
			total = 0
		}
		if !cmdOptions.IsOptionSet("compact") {
			//colors.DisplayOk(" [ " + startTime.Format("15:04") + " -> " + endTime.Format("15:04") + " ] " + duration.String() + " : " + item.Summary)
		}
	}

	for i := 0; i < len(datesEvent); i++ {
		datesEvent[i] = time.Date(datesEvent[i].Year(), datesEvent[i].Month(), datesEvent[i].Day(), 0, 0, 0, 0, time.Local)
		s := fmt.Sprintf("%f", durationWork[i])

		colors.DisplayInfo(datesEvent[i].String() + "   ->  " + s)
	}

	graph := chart.Chart{
		Series: []chart.Series{
			chart.TimeSeries{
				XValues: datesEvent,
				YValues: durationWork,
			},
		},
	}

	f, _ := os.Create("output.png")
	defer f.Close()
	graph.Render(chart.PNG, f)
	return nil
}
