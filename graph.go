package studio_statistics

import (
	"encoding/csv"
	"fmt"
	"github.com/wcharczuk/go-chart"
	"io"
	"regexp"
	"strconv"
	"time"
)

var alphanumeric = regexp.MustCompile("(?m)[^a-zA-Z0-9, \\-:.]+")

func MakeGraph(r io.Reader, w io.Writer) error {
	rows, err := csv.NewReader(r).ReadAll()
	if err != nil {
		return err
	}
	fmt.Println("Getting rid of identical rows...")
	before := len(rows)
	rows = eliminateIdenticalTimestamps(rows)
	after := len(rows)
	fmt.Println("Eliminated", before-after, "rows!")
	fmt.Println("Parsing dataset")
	times := make([]time.Time, 0, len(rows))
	opens := make([]float64, 0, len(rows))
	switchStates := make([]float64, 0, len(rows))
	motions := make([]float64, 0, len(rows))
	for i := 0; i < len(rows); i++ {
		for j := 0; j < len(rows[i]); j++ {
			rows[i][j] = alphanumeric.ReplaceAllString(rows[i][j], "")
		}
		t, err := parseTime(rows[i][0])
		if t.Add(time.Hour * 24 * 7).Before(time.Now()) {
			continue
		}
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
	fmt.Println("Done parsing data")
	width := int(times[len(times)-1].Sub(times[0]).Seconds() / 30.0)
	graph := chart.Chart{
		Title:  "Design Studio Statistics",
		Width:  width,
		Height: 600,
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
	fmt.Println("Rendering graph")
	err = graph.Render(chart.PNG, w)
	if err != nil {
		return err
	}

	return nil
}

func btof(v bool) float64 {
	if v {
		return 1.0
	}
	return 0.0
}

func eliminateIdenticalTimestamps(rows [][]string) (newRows [][]string) {
	if rows == nil || len(rows) < 3 {
		newRows = rows
		return
	}
	for i := 0; i < len(rows); i++ {
		lastIdenticalIndex := i
		for j := i; j < len(rows); j++ {
			identical := true
			for k := 1; k < len(rows[i]); k++ {
				if rows[i][k] != rows[j][k] {
					identical = false
					break
				}
			}
			if identical {
				lastIdenticalIndex = j
			} else {
				break
			}
		}
		if lastIdenticalIndex != i {
			newRows = append(newRows, rows[i], rows[lastIdenticalIndex])
			i = lastIdenticalIndex
		} else {
			newRows = append(newRows, rows[i])
		}
	}
	return
}

func parseTime(timeString string) (t time.Time, err error) {
	t, err= time.Parse(time.RFC822, timeString)
	if err != nil {
		t, err = time.Parse(time.RFC3339Nano, timeString)
	}
	return
}