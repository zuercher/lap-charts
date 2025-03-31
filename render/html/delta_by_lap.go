package html

import (
	"embed"
	"fmt"
	"html/template"
	"math"

	"zuercher.us/lapcharts/mode"
	"zuercher.us/lapcharts/render"
	"zuercher.us/lapcharts/source"
	"zuercher.us/lapcharts/util/iox"
	"zuercher.us/lapcharts/util/timex"
)

//go:embed delta_by_lap.html.tmpl
var deltaByLapTemplate embed.FS

const numYTicks = 5

type DeltaByLapOptions struct{}

func (o *DeltaByLapOptions) RenderMode() mode.RenderMode {
	return mode.HtmlDeltaByLap
}
func (o *DeltaByLapOptions) ConfigureFlags() {
}
func (o *DeltaByLapOptions) Validate() error {
	return nil
}
func (o *DeltaByLapOptions) Renderer(source source.Source) (render.Renderer, error) {
	return &deltaByLapRenderer{Source: source}, nil
}

type deltaByLapRenderer struct {
	source.Source

	byCar  []deltaByLapCar
	yTicks []timex.Duration
}

func (r *deltaByLapRenderer) Render(output iox.Writer) error {
	tmpl, err :=
		template.New("delta_by_lap.html.tmpl").
			Funcs(templateFuncs()).
			ParseFS(deltaByLapTemplate, "delta_by_lap.html.tmpl")
	if err != nil {
		return err
	}

	return tmpl.Execute(output, r)
}

func (r *deltaByLapRenderer) Cars() []deltaByLapCar {
	if r.byCar == nil {
		// Relative to the first car.
		for idx := range r.Source.NumCars() - 1 {
			first := r.Source.Laps(0)
			firstCumulative, _ := first.StartOffset()
			this := r.Source.Laps(idx + 1)
			thisCumulative, _ := this.StartOffset()

			numLaps := min(this.NumLaps(), first.NumLaps())
			deltas := make([]timex.Duration, numLaps)

			for idx := range numLaps {
				firstLap := first.Lap(idx + 1)
				thisLap := this.Lap(idx + 1)

				thisCumulative += thisLap.Time
				firstCumulative += firstLap.Time

				deltas[idx] = thisCumulative - firstCumulative
			}
			r.byCar = append(r.byCar, deltaByLapCar{deltas: deltas})
		}
	}
	return r.byCar
}

func (r *deltaByLapRenderer) YTicks() []timex.Duration {
	if r.yTicks == nil {
		// Get the min and max values from all cars
		minDelta := timex.Duration(0)
		maxDelta := timex.Duration(0)
		for _, car := range r.Cars() {
			for _, delta := range car.Deltas() {
				minDelta = min(minDelta, delta)
				maxDelta = max(maxDelta, delta)
			}
		}

		if maxDelta < 0 {
			maxDelta = 0
		}

		if minDelta > 0 {
			minDelta = 0
		}

		if minDelta == maxDelta {
			r.yTicks = []timex.Duration{maxDelta}
			return r.yTicks
		}

		var tickSize timex.Duration
		if minDelta < 0 && maxDelta > 0 {
			// Compute tick size based on whether we're more negative or more positive.
			if 0-minDelta > maxDelta {
				tickSize = computeTickSize(-minDelta, numYTicks)
			} else {
				tickSize = computeTickSize(maxDelta, numYTicks)
			}
		} else if maxDelta == 0 {
			// All are negative.
			tickSize = computeTickSize(-minDelta, numYTicks)
		} else {
			// All are postive.
			tickSize = computeTickSize(maxDelta, numYTicks)
		}

		negativeTicks := 0
		for minDelta <= timex.Duration(negativeTicks)*tickSize {
			negativeTicks--
		}
		positiveTicks := 0
		for maxDelta >= timex.Duration(positiveTicks)*tickSize {
			positiveTicks++
		}

		for tick := negativeTicks; tick <= positiveTicks; tick++ {
			r.yTicks = append(r.yTicks, timex.Duration(tick)*tickSize)
		}
	}
	return r.yTicks
}

func (r *deltaByLapRenderer) YTicksFormatted() []string {
	ticks := []string{}
	for _, tick := range r.YTicks() {
		ticks = append(ticks, tick.String())
	}
	return ticks
}

func computeTickSize(m timex.Duration, numTicks int) timex.Duration {
	m = m / timex.Duration(numTicks)

	if m == 0 {
		return timex.Second
	}

	var multiplier timex.Duration
	buckets := []float64{0.1, 0.2, 0.25, 0.3, 0.4, 0.5, 1.0}
	if m < timex.Second {
		multiplier = timex.Millisecond
	} else if m < timex.Minute {
		// Compute in seconds
		multiplier = timex.Second
		buckets = []float64{0.05, 0.1, 0.2, 0.3, 0.5, 1.0}
	} else {
		// Compute in minutes
		multiplier = timex.Minute
	}

	rawTick := float64(m / multiplier)
	x := 0
	v := 0.0
	for {
		v = rawTick / math.Pow10(x)
		if v < 1.0 {
			break
		}
		x++
		if x >= 7 {
			// Give up and pick something.
			return m
		}
	}

	for _, bucket := range buckets {
		if v <= bucket {
			v = bucket
			break
		}
	}

	return timex.Duration(math.Round(v*math.Pow10(x))) * multiplier
}

type deltaByLapCar struct {
	deltas []timex.Duration
}

func (c *deltaByLapCar) LapNumbers() []int {
	lapNumbers := make([]int, len(c.deltas))
	for idx := range c.deltas {
		lapNumbers[idx] = idx + 1
	}
	return lapNumbers
}

func (c *deltaByLapCar) Deltas() []timex.Duration {
	return c.deltas
}

func (c *deltaByLapCar) HoverText() []string {
	hovers := make([]string, len(c.deltas))
	for i, t := range c.deltas {
		hovers[i] = fmt.Sprintf("Lap %d: %s", i+1, t)
	}
	return hovers
}
