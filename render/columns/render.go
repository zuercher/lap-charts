package columns

import (
	"fmt"
	"strings"

	"zuercher.us/lapcharts/source"
	"zuercher.us/lapcharts/util/iox"
	"zuercher.us/lapcharts/util/timex"
)

type columnRenderer struct {
	source source.Source
}

func (r *columnRenderer) Render(w iox.Writer) error {
	rows := [][]string{}
	cols := []source.Laps{}
	cumulativeTimes := []timex.Duration{}

	for car := range r.source.NumCars() {
		laps := r.source.Laps(car)
		cols = append(cols, laps)

		offset, _ := laps.StartOffset()
		cumulativeTimes = append(cumulativeTimes, offset)
	}
	lap := 1
	maxes := make([]int, len(cols)+1)
	for {
		row := []string{}

		lapDesc := fmt.Sprintf("Lap %d", lap)
		maxes[0] = max(maxes[0], len(lapDesc))
		row = append(row, lapDesc)

		done := 0
		for idx, laps := range cols {
			if num := laps.NumLaps(); lap <= num {
				t := laps.Lap(lap).Time
				cumulativeTimes[idx] += t
				tStr := t.String()
				row = append(row, tStr)
				if lap == num {
					done++
				}

				maxes[idx+1] = max(maxes[idx+1], len(tStr))
			} else {
				row = append(row, "-")
				done++
			}
		}
		rows = append(rows, row)
		if done == len(cols) {
			break
		}
		lap++
	}

	for idx, t := range cumulativeTimes {
		maxes[idx+1] = max(maxes[idx+1], len(t.String()))
	}

	for _, row := range rows {
		for idx, col := range row {
			if idx > 0 {
				w.Print("  ")
			}
			w.Printf("%*s", maxes[idx], col)
		}
		w.Println()
	}

	w.Print(strings.Repeat(" ", maxes[0]))
	for idx, t := range cumulativeTimes {
		w.Print("  ")
		w.Printf("%*s", maxes[idx+1], t.String())
	}
	w.Println()

	return nil
}
