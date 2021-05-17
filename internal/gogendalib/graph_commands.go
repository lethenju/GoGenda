package gogendalib

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	cmdOptions "github.com/lethenju/gogenda/internal/cmd_options"
	"github.com/lethenju/gogenda/internal/configuration"
	"github.com/lethenju/gogenda/internal/utilities"
	"github.com/lethenju/gogenda/pkg/colors"
	api "github.com/lethenju/gogenda/pkg/google_agenda_api"
	"google.golang.org/api/calendar/v3"
)

func RenderGraphCompleteness(items []*calendar.Event) *charts.Bar {

	var durationTotal []float64
	for i := 0; i < 8; i++ {
		durationTotal = append(durationTotal, 0)
	}
	for _, item := range items {
		startTime, _ := time.Parse(time.RFC3339, item.Start.DateTime)
		endTime, _ := time.Parse(time.RFC3339, item.End.DateTime)
		duration := endTime.Sub(startTime)

		durationTotal[startTime.Weekday()] += duration.Hours() / float64(len(items))
	}

	// create a new bar instance
	bar := charts.NewBar()
	// set some global options like Title/Legend/ToolTip or anything else
	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title: "Completeness since " + items[0].Start.Date,
	}), charts.WithToolboxOpts(opts.Toolbox{Show: true}),
		charts.WithLegendOpts(opts.Legend{Right: "80%"}))

	itemsTotal := make([]opts.BarData, 0)
	for i := 1; i < 7; i++ {
		itemsTotal = append(itemsTotal, opts.BarData{Value: durationTotal[i]})
	}
	// Sunday at last
	itemsTotal = append(itemsTotal, opts.BarData{Value: durationTotal[0]})

	// Put data into instance
	bar.SetXAxis([]string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}).
		AddSeries("Total", itemsTotal)
	return bar
}

func RenderGraphWorkVsLastWeek(items []*calendar.Event) *charts.Bar {

	// We have to put a default 0 value that is not the zero value of item.ColorId (which is ""
	var durationWork []float64
	// last week
	var durationWorkLastWeek []float64
	// The week before
	var durationWorkLastWeek2 []float64
	for i := 0; i < 8; i++ {
		durationWork = append(durationWork, 0)
		durationWorkLastWeek = append(durationWorkLastWeek, 0)
		durationWorkLastWeek2 = append(durationWorkLastWeek2, 0)
	}
	nowYear, nowWeek := time.Now().ISOWeek()

	for _, item := range items {
		startTime, _ := time.Parse(time.RFC3339, item.Start.DateTime)
		endTime, _ := time.Parse(time.RFC3339, item.End.DateTime)
		duration := endTime.Sub(startTime)

		// retrieve category
		colorName, _ := api.GetColorNameFromColorID(item.ColorId)
		category := configuration.GetNameFromColor(colorName)
		itemYear, itemWeek := startTime.ISOWeek()

		if category == "WORK" {
			if nowYear == itemYear && nowWeek == itemWeek {
				durationWork[startTime.Weekday()] += duration.Hours()
			} else if nowYear == itemYear && nowWeek == itemWeek+1 {
				// Last week (TODO not working for last / first week of the year)
				durationWorkLastWeek[startTime.Weekday()] += duration.Hours()
			} else if nowYear == itemYear && nowWeek == itemWeek+2 {
				// the week before
				durationWorkLastWeek2[startTime.Weekday()] += duration.Hours()
			}
		}
	}

	// create a new bar instance
	bar := charts.NewBar()
	// set some global options like Title/Legend/ToolTip or anything else
	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title: "Worktime this week vs last weeks",
	}), charts.WithToolboxOpts(opts.Toolbox{Show: true}),
		charts.WithLegendOpts(opts.Legend{Right: "80%"}))

	itemsWork := make([]opts.BarData, 0)
	itemsWorkLastWeek := make([]opts.BarData, 0)
	itemsWorkLastWeek2 := make([]opts.BarData, 0)
	for i := 1; i < 7; i++ {
		itemsWork = append(itemsWork, opts.BarData{Value: durationWork[i]})
		itemsWorkLastWeek = append(itemsWorkLastWeek, opts.BarData{Value: durationWorkLastWeek[i]})
		itemsWorkLastWeek2 = append(itemsWorkLastWeek2, opts.BarData{Value: durationWorkLastWeek2[i]})
	}
	// Sunday at last
	itemsWork = append(itemsWork, opts.BarData{Value: durationWork[0]})
	itemsWorkLastWeek = append(itemsWorkLastWeek, opts.BarData{Value: durationWorkLastWeek[0]})
	itemsWorkLastWeek2 = append(itemsWorkLastWeek2, opts.BarData{Value: durationWorkLastWeek2[0]})

	// Put data into instance
	bar.SetXAxis([]string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}).
		AddSeries("Work the week before", itemsWorkLastWeek2).
		AddSeries("Work last week", itemsWorkLastWeek).
		AddSeries("Work", itemsWork)
	//AddSeries("Projets", itemsProjects).
	return bar
}

func RenderGraphWorkVsPlay(items []*calendar.Event) *charts.Bar {

	// We have to put a default 0 value that is not the zero value of item.ColorId (which is ""
	var durationWork []float64
	var durationFun []float64
	for i := 0; i < 8; i++ {
		durationWork = append(durationWork, 0)
		durationFun = append(durationFun, 0)
	}
	for _, item := range items {
		startTime, _ := time.Parse(time.RFC3339, item.Start.DateTime)
		endTime, _ := time.Parse(time.RFC3339, item.End.DateTime)
		duration := endTime.Sub(startTime)

		// retrieve category
		colorName, _ := api.GetColorNameFromColorID(item.ColorId)
		category := configuration.GetNameFromColor(colorName)

		if category == "WORK" {
			durationWork[startTime.Weekday()] += duration.Hours()
		} else if category == "FUN" {
			durationFun[startTime.Weekday()] += duration.Hours()
		} else if category == "PROJECT" {
			durationWork[startTime.Weekday()] += duration.Hours()
		}

		if !cmdOptions.IsOptionSet("compact") {
			//colors.DisplayOk(" [ " + startTime.Format(time.RFC3339) + " -> " + endTime.Format("15:04") + " ] " + duration.String() + " : " + item.Summary)
		}
	}

	// create a new bar instance
	bar := charts.NewBar()
	// set some global options like Title/Legend/ToolTip or anything else
	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    "Worktime vs Playtime since " + items[0].Start.Date,
		Subtitle: "In blue worktime, in green playtime",
	}), charts.WithToolboxOpts(opts.Toolbox{Show: true}),
		charts.WithLegendOpts(opts.Legend{Right: "80%"}))

	itemsWork := make([]opts.BarData, 0)
	itemsFun := make([]opts.BarData, 0)
	for i := 1; i < 7; i++ {
		itemsWork = append(itemsWork, opts.BarData{Value: durationWork[i]})
		itemsFun = append(itemsFun, opts.BarData{Value: durationFun[i]})
	}
	// Sunday at last
	itemsWork = append(itemsWork, opts.BarData{Value: durationWork[0]})
	itemsFun = append(itemsFun, opts.BarData{Value: durationFun[0]})

	// Put data into instance
	bar.SetXAxis([]string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}).
		AddSeries("Work", itemsWork).
		//AddSeries("Projets", itemsProjects).
		AddSeries("Fun", itemsFun)
	return bar
}

func RenderGraphWork(items []*calendar.Event) *charts.Line {
	line := charts.NewLine()

	x := make([]string, 0)
	//y := make([]opts.LineData, 0)
	y2 := make([]opts.LineData, 0)
	y3 := make([]opts.LineData, 0)
	actualDate, _ := time.Parse(time.RFC3339, items[0].Start.DateTime)
	totalDuration := time.Duration(0)
	totalDurationWork := time.Duration(0)
	totalDurationFun := time.Duration(0)

	for i := 0; i < len(items); i++ {
		startTime, _ := time.Parse(time.RFC3339, items[i].Start.DateTime)
		endTime, _ := time.Parse(time.RFC3339, items[i].End.DateTime)
		duration := endTime.Sub(startTime)

		// Same date
		_, week := startTime.ISOWeek()
		_, weekAct := actualDate.ISOWeek()
		if week == weekAct {
			totalDuration += duration

			// retrieve category
			colorName, _ := api.GetColorNameFromColorID(items[i].ColorId)
			category := configuration.GetNameFromColor(colorName)

			if category == "WORK" || category == "PROJECT" {
				totalDurationWork += duration
			} else if category == "FUN" {
				totalDurationFun += duration
			}
		} else {
			if totalDuration != 0 {
				x = append(x, startTime.Format(time.RFC1123))
				//y = append(y, opts.LineData{Value: totalDuration.Hours()})
				y2 = append(y2, opts.LineData{Value: (totalDurationWork.Hours() / totalDuration.Hours())})
				y3 = append(y3, opts.LineData{Value: (totalDurationFun.Hours() / totalDuration.Hours())})
			}

			totalDuration = time.Duration(0)
			totalDurationWork = time.Duration(0)
			totalDurationFun = time.Duration(0)
			// different date, add data
			actualDate = startTime
		}
	}

	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Event type evolution",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			SplitNumber: 20,
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Scale: true,
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "inside",
			Start:      50,
			End:        100,
			XAxisIndex: []int{0},
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "slider",
			Start:      50,
			End:        100,
			XAxisIndex: []int{0},
		}),
	)

	line.SetXAxis(x).AddSeries("LineFun", y3).AddSeries("LineWork", y2).SetSeriesOptions(charts.WithLineChartOpts(
		opts.LineChart{
			Smooth: true,
		}),
	)

	return line
}
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

	// If there are too much time between dates, need to ask several times : ask for every 30 days
	nbAsks := (int(end.Sub(begin).Hours()/24) / 30) + 1
	colors.DisplayInfo("NbAsks " + strconv.Itoa(nbAsks))
	var items []*calendar.Event

	if nbAsks == 1 {
		events, err := api.GetActivitiesBetweenDates(
			begin.Format(time.RFC3339),
			end.Format(time.RFC3339), srv)
		if err != nil {
			return err
		}

		items = append(items, events.Items...)
	} else {

		for i := 1; i < nbAsks; i++ {
			colors.DisplayInfo("Asking for 30 days (i=" + strconv.Itoa(i) + ") - [" + begin.Add(time.Hour*time.Duration(24*30*(i-1))).Format(time.RFC3339) + "] -> [" + begin.Add(time.Hour*time.Duration(24*30*i)).Format(time.RFC3339) + "]")

			//Asking for 30 days
			events, err := api.GetActivitiesBetweenDates(
				begin.Add(time.Hour*time.Duration(24*30*(i-1))).Format(time.RFC3339),
				begin.Add(time.Hour*time.Duration(24*30*i)).Format(time.RFC3339), srv)
			if err != nil {
				return err
			}
			items = append(items, events.Items...)
		}
		// Asking for the remaining days
		colors.DisplayInfo("Asking for the remaining days - [" + begin.Add(time.Hour*time.Duration(24*120*(nbAsks-1))).Format(time.RFC3339) + "] -> [" + end.Format(time.RFC3339) + "]")
		events, err := api.GetActivitiesBetweenDates(
			begin.Add(time.Hour*time.Duration(24*120*(nbAsks-1))).Format(time.RFC3339),
			end.Format(time.RFC3339), srv)
		if err != nil {
			return err
		}
		items = append(items, events.Items...)
	}

	page := components.NewPage()
	page.Layout = components.PageFlexLayout
	page.AddCharts(
		RenderGraphWorkVsPlay(items),
		RenderGraphWorkVsLastWeek(items),
		RenderGraphCompleteness(items),
		RenderGraphWork(items),
	)
	// Where the magic happens
	f, _ := os.Create("page.html")
	page.Render(f)
	return nil
}
