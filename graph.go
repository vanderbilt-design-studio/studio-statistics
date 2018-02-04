package studio_statistics

import (
	"encoding/csv"
	"github.com/wcharczuk/go-chart"
	"regexp"
	"strconv"
	"time"
	"io"
)

var alphanumeric = regexp.MustCompile("(?m)[^a-zA-Z0-9, :]+")

func MakeGraph(r io.Reader, w io.Writer) (error) {
	rows, err := csv.NewReader(r).ReadAll()
	if err != nil {
		return err
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
			return err
		}
		switchState, err := strconv.ParseInt(rows[i][2], 10, 32)
		if err != nil {
			return err
		}
		motion, err := strconv.ParseBool(rows[i][3])
		if err != nil {
			return err
		}
		times = append(times, t)
		opens = append(opens, btof(open))
		switchStates = append(switchStates, float64(switchState))
		motions = append(motions, btof(motion))
	}
	graph := chart.Chart{
		Title:  "Design Studio Statistics",
		Width:  1920 * 10,
		Height: 800,
		XAxis: chart.XAxis{
			Name:           "Time",
			NameStyle:      chart.StyleShow(),
			Style:          chart.StyleShow(),
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
	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}
	err = graph.Render(chart.PNG, w)
	if err != nil {
		return err
	}

	return nil

	//fcsvout, err := os.Create("graph.csv")
	//if err != nil {
	//	panic(err)
	//}
	//rows = make([][]string, len(times))
	//for i := range times {
	//	rows[i] = []string{
	//		times[i].Format(time.RFC3339Nano),
	//		fmt.Sprintf("%v", opens[i]),
	//		fmt.Sprintf("%v", switchStates[i]),
	//		fmt.Sprintf("%v", motions[i]),
	//	}
	//}
	//w := csv.NewWriter(fcsvout)
	//w.WriteAll(rows)
	//w.Flush()
	//fcsvout.Close()
}

func btof(v bool) float64 {
	if v {
		return 1.0
	}
	return 0.0
}
