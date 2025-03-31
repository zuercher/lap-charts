package html

import (
	"embed"
	"fmt"
	"html/template"

	"zuercher.us/lapcharts/mode"
	"zuercher.us/lapcharts/render"
	"zuercher.us/lapcharts/source"
	"zuercher.us/lapcharts/util/iox"
	"zuercher.us/lapcharts/util/timex"
)

//go:embed laps_by_time.html.tmpl
var lapsByTimeTemplate embed.FS

type LapsByTimeOptions struct {
}

func (o *LapsByTimeOptions) RenderMode() mode.RenderMode {
	return mode.HtmlLapsByTime
}
func (o *LapsByTimeOptions) ConfigureFlags() {
}
func (o *LapsByTimeOptions) Validate() error {
	return nil
}
func (o *LapsByTimeOptions) Renderer(source source.Source) (render.Renderer, error) {
	return &lapsByTimeRenderer{Source: source}, nil
}

type lapsByTimeRenderer struct {
	source.Source

	byCar []lapsByTimeCar
	ticks []timex.Duration
}

func (r *lapsByTimeRenderer) Render(output iox.Writer) error {
	tmpl, err :=
		template.New("laps_by_time.html.tmpl").
			Funcs(templateFuncs()).
			ParseFS(lapsByTimeTemplate, "laps_by_time.html.tmpl")
	if err != nil {
		return err
	}

	return tmpl.Execute(output, r)
}

func (r *lapsByTimeRenderer) Cars() []lapsByTimeCar {
	if r.byCar == nil {
		cars := []lapsByTimeCar{}
		for idx := range r.Source.NumCars() {
			laps := r.Source.Laps(idx)

			startOffset, _ := laps.StartOffset()

			lapTimes := make([]timex.Duration, laps.NumLaps())
			for lap := range laps.NumLaps() {
				lapTimes[lap] = laps.Lap(lap + 1).Time
			}

			lbtCar := lapsByTimeCar{
				startOffset: startOffset,
				lapTimes:    lapTimes,
			}
			cars = append(cars, lbtCar)
		}
		r.byCar = cars
	}
	return r.byCar
}

func (c *lapsByTimeRenderer) Ticks() []timex.Duration {
	if c.ticks == nil {
		last := timex.Duration(0)
		for _, car := range c.Cars() {
			carLast := car.LapTimes()[len(car.LapTimes())-1]
			if carLast > last {
				last = carLast
			}
		}

		ticks := []timex.Duration{}
		increment := timex.Duration(5 * timex.Minute)
		for tick := timex.Duration(0); tick < last; tick += increment {
			ticks = append(ticks, tick)
		}
		c.ticks = append(ticks, ticks[len(ticks)-1]+increment)
	}
	return c.ticks
}

func (c *lapsByTimeRenderer) TicksFormatted() []string {
	ticks := make([]string, len(c.Ticks()))
	for i, t := range c.Ticks() {
		ticks[i] = t.String()
	}
	return ticks
}

type lapsByTimeCar struct {
	startOffset          timex.Duration
	lapTimes             []timex.Duration
	cumulativeLapTimesMs []timex.Duration
}

// LapTimes returns cumulative lap times in milliseconds.
func (c *lapsByTimeCar) LapTimes() []timex.Duration {
	if c.cumulativeLapTimesMs == nil {
		lapTimes := make([]timex.Duration, len(c.lapTimes))
		cumulative := c.startOffset
		for i, t := range c.lapTimes {
			cumulative += t
			lapTimes[i] = cumulative
		}
		c.cumulativeLapTimesMs = lapTimes
	}
	return c.cumulativeLapTimesMs
}

func (c *lapsByTimeCar) LapTimesFormatted() []string {
	lapTimes := make([]string, len(c.lapTimes))
	cumulative := c.startOffset
	for i, t := range c.lapTimes {
		cumulative += t
		lapTimes[i] = cumulative.String()
	}
	return lapTimes
}

func (c *lapsByTimeCar) LapNumbers() []int {
	lapNumbers := make([]int, len(c.lapTimes))
	for i := range c.lapTimes {
		lapNumbers[i] = i + 1
	}
	return lapNumbers
}

func (c *lapsByTimeCar) HoverText() []string {
	hovers := make([]string, len(c.lapTimes))
	for i, t := range c.lapTimes {
		hovers[i] = fmt.Sprintf("Lap %d: %s", i+1, t)
	}
	return hovers
}
