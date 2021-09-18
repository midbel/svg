package main

import (
	"os"
	"sort"
	"time"

	"github.com/midbel/svg"
	"github.com/midbel/svg/chart"
)

func main() {
	var (
		w  = os.Stdout
		c  chart.GanttChart
		i1 = chart.Interval{
			Title:  "task-1",
			Starts: time.Date(2021, 9, 1, 10, 0, 0, 0, time.UTC),
			Ends:   time.Date(2021, 9, 4, 17, 0, 0, 0, time.UTC),
			Fill:   svg.NewFill("olive"),
		}
		i2 = chart.Interval{
			Title:  "task-2",
			Starts: time.Date(2021, 9, 7, 10, 0, 0, 0, time.UTC),
			Ends:   time.Date(2021, 9, 10, 17, 0, 0, 0, time.UTC),
			Fill:   svg.NewFill("orchid"),
		}
		i3 = chart.Interval{
			Title:  "task-3",
			Starts: time.Date(2021, 9, 1, 14, 0, 0, 0, time.UTC),
			Ends:   time.Date(2021, 9, 10, 12, 0, 0, 0, time.UTC),
			Fill:   svg.NewFill("salmon"),
		}
		i4 = chart.Interval{
			Title:  "task-4",
			Starts: time.Date(2021, 9, 2, 14, 0, 0, 0, time.UTC),
			Ends:   time.Date(2021, 9, 2, 17, 0, 0, 0, time.UTC),
			Fill:   svg.NewFill("steelblue"),
		}
		i5 = chart.Interval{
			Title:  "task-5",
			Starts: time.Date(2021, 9, 3, 10, 0, 0, 0, time.UTC),
			Ends:   time.Date(2021, 9, 5, 17, 0, 0, 0, time.UTC),
			Fill:   svg.NewFill("steelblue"),
		}
		i6 = chart.Interval{
			Title:  "task-6",
			Starts: time.Date(2021, 9, 6, 8, 0, 0, 0, time.UTC),
			Ends:   time.Date(2021, 9, 10, 20, 0, 0, 0, time.UTC),
			Fill:   svg.NewFill("steelblue"),
		}
		i7 = chart.Interval{
			Title:  "task-7",
			Starts: time.Date(2021, 9, 2, 10, 0, 0, 0, time.UTC),
			Ends:   time.Date(2021, 9, 4, 17, 0, 0, 0, time.UTC),
			Fill:   svg.NewFill("olive"),
		}
		i8 = chart.Interval{
			Title:  "task-8",
			Starts: time.Date(2021, 9, 5, 10, 0, 0, 0, time.UTC),
			Ends:   time.Date(2021, 9, 9, 17, 0, 0, 0, time.UTC),
			Fill:   svg.NewFill("olive"),
		}
	)
	i1.Sub = append(i1.Sub, i8, i6)
	i4.Sub = append(i4.Sub, i5, i2)
	c.Width = 1920
	c.Height = 640
	c.Padding = chart.Padding{
		Left:   80,
		Right:  40,
		Top:    20,
		Bottom: 40,
	}
	c.Axis.Left = chart.CreateLabelAxis()
	c.Axis.Bottom = chart.CreateTimeAxis(chart.WithTicks(10))
	xs := []chart.Interval{i1, i3, i4, i7}
	sort.Slice(xs, func(i, j int) bool {
		return xs[i].Starts.Before(xs[j].Starts)
	})

	c.Render(w, xs)
}
