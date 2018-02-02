package main

import (
	"encoding/csv"
	"fmt"
	"github.com/wcharczuk/go-chart"
	"os"
	"regexp"
	"strconv"
	"time"
)

var alphanumeric = regexp.MustCompile("(?m)[^a-zA-Z0-9, :]+")

func main() {
	file, err := os.Open("activity.log")
	if err != nil {
		panic("Failed to open file")
	}
	rows, err := csv.NewReader(file).ReadAll()
	if err != nil {
		panic("Failed to read from file")
	}
	times := make([]time.Time, 0, len(rows))
	opens := make([]float64, 0, len(rows))
	switchStates := make([]float64, 0, len(rows))
	motions := make([]float64, 0, len(rows))
	for i := 0; i < len(rows); i++ {
		for j := 0; j < len(rows[i]); j++ {
			rows[i][j] = alphanumeric.ReplaceAllString(rows[i][j], "")
		}
		t, err := time.Parse("02 Jan 06 15:04 MST", rows[i][0])
		if err != nil {
			panic(err)
		}
		/*
		var topen, tclose = time.Date(t.Year(), t.Month(), t.Day(), 12, 00, 00, 0, time.Local), time.Date(t.Year(), t.Month(), t.Day(), 22, 00, 00, 0, time.Local)
		if t.Before(topen) || t.After(tclose) {
			continue
		}
		*/
		open, err := strconv.ParseBool(rows[i][1])
		if err != nil {
			panic("Failed to parse open state")
		}
		switchState, err := strconv.ParseInt(rows[i][2], 10, 32)
		if err != nil {
			panic("Failed to parse switchState")
		}
		motion, err := strconv.ParseBool(rows[i][3])
		if err != nil {
			panic("Failed to parse if motion")
		}
		times = append(times, t)
		opens = append(opens, Btof(open))
		switchStates = append(switchStates, float64(switchState))
		motions = append(motions, Btof(motion))
	}
	graph := chart.Chart{
		Title:  "Design Studio Statistics",
		Width:  1920*10,
		Height: 800,
		XAxis: chart.XAxis{
			Name:      "Time",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
			ValueFormatter: chart.TimeValueFormatterWithFormat("Mon Jan _2 3:04PM"),
		},
		YAxis: chart.YAxis{
			Name:      "Value",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
		},
		Series: []chart.Series{
			chart.TimeSeries{
				Name:    "Time Open",
				XValues: times,
				YValues: opens,
			},
			chart.TimeSeries{
				Name:    "Switch States (0 = normal, 1 = forced open, 2 = forced closed)",
				XValues: times,
				YValues: switchStates,
			},
			chart.TimeSeries{
				Name:    "Motion Detected",
				XValues: times,
				YValues: motions,
			},
		},
	}
	fout, err := os.Create("graph.png")
	if err != nil {
		panic("Failed to open file for output")
	}
	graph.Elements = []chart.Renderable{
		chart.LegendLeft(&graph),
	}
	err = graph.Render(chart.PNG, fout)
	if err != nil {
		panic(err)
	}
	fmt.Println("Done!")
}

func Btof(v bool) float64 {
	if v {
		return 1.0
	}
	return 0.0
}
