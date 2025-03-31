package source

import (
	"zuercher.us/lapcharts/util/ptr"
	"zuercher.us/lapcharts/util/timex"
)

type Source interface {
	NumCars() int
	Laps(car int) Laps
}

type Laps struct {
	description string
	lap         []Lap
	offset      *timex.Duration
}

func NewLaps(description string) Laps {
	return Laps{
		description: description,
	}
}

func (l *Laps) Description() string { return l.description }
func (l *Laps) NumLaps() int        { return len(l.lap) }
func (l *Laps) Lap(lap int) Lap     { return l.lap[lap-1] }
func (l *Laps) AppendLap(lap Lap)   { l.lap = append(l.lap, lap) }
func (l *Laps) LastLap() (Lap, bool) {
	if len(l.lap) == 0 {
		return Lap{}, false
	}
	return l.lap[len(l.lap)-1], true
}
func (l *Laps) StartOffset() (timex.Duration, bool)  { return ptr.Get(l.offset), l.offset != nil }
func (l *Laps) SetStartOffset(offset timex.Duration) { l.offset = &offset }

type Lap struct {
	Num      int
	Time     timex.Duration
	Position int
	SpeedKPH float64
}
