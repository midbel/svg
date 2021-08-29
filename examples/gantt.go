package main

import (
	"os"
	"time"

	"github.com/midbel/svg"
	"github.com/midbel/svg/chart"
)

func main() {
	var (
		w  = os.Stdout
		xs []chart.GanttSerie
		c  chart.GanttChart
		i1 = chart.Interval{
			Title:  "task-1",
			Starts: time.Date(2021, 9, 1, 10, 0, 0, 0, time.UTC),
			Ends:   time.Date(2021, 9, 4, 17, 0, 0, 0, time.UTC),
		}
		i2 = chart.Interval{
			Title:  "task-2",
			Starts: time.Date(2021, 9, 7, 10, 0, 0, 0, time.UTC),
			Ends:   time.Date(2021, 9, 10, 17, 0, 0, 0, time.UTC),
		}
		i3 = chart.Interval{
			Title:  "task-3",
			Starts: time.Date(2021, 9, 1, 14, 0, 0, 0, time.UTC),
			Ends:   time.Date(2021, 9, 6, 17, 0, 0, 0, time.UTC),
		}
		i4 = chart.Interval{
			Title:  "task-4",
			Starts: time.Date(2021, 9, 2, 14, 0, 0, 0, time.UTC),
			Ends:   time.Date(2021, 9, 2, 17, 0, 0, 0, time.UTC),
		}
		i5 = chart.Interval{
			Title:  "task-5",
			Starts: time.Date(2021, 9, 3, 10, 0, 0, 0, time.UTC),
			Ends:   time.Date(2021, 9, 5, 17, 0, 0, 0, time.UTC),
		}
		i6 = chart.Interval{
			Title:  "task-6",
			Starts: time.Date(2021, 9, 6, 8, 0, 0, 0, time.UTC),
			Ends:   time.Date(2021, 9, 10, 12, 0, 0, 0, time.UTC),
		}
		i7 = chart.Interval{
			Title:  "task-7",
			Starts: time.Date(2021, 9, 3, 15, 0, 0, 0, time.UTC),
			Ends:   time.Date(2021, 9, 4, 17, 0, 0, 0, time.UTC),
		}
		i8 = chart.Interval{
			Title:  "task-8",
			Starts: time.Date(2021, 9, 5, 10, 0, 0, 0, time.UTC),
			Ends:   time.Date(2021, 9, 9, 17, 0, 0, 0, time.UTC),
		}
	)
	c.Width = 960
	c.Height = 480
	c.Padding = chart.CreatePadding(20, 20)
	xs = append(xs, chart.NewGanttSerie("serie-1"))
	xs[0].Fill = svg.NewFill("salmon")
	xs[0].Append(i1)
	xs[0].Append(i2)
	xs = append(xs, chart.NewGanttSerie("serie-2"))
	xs[1].Fill = svg.NewFill("orchid")
	xs[1].Append(i3)
	xs = append(xs, chart.NewGanttSerie("serie-3"))
	xs[2].Fill = svg.NewFill("olive")
	xs[2].Append(i4)
	xs[2].Append(i5)
	xs[2].Append(i6)
	xs = append(xs, chart.NewGanttSerie("serie-4"))
	xs[3].Fill = svg.NewFill("steelblue")
	xs[3].Append(i7)
	xs[3].Append(i8)
	c.Render(w, xs)
}