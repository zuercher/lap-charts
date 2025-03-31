package timex

import (
	"fmt"
	"time"
)

// Duration in milliseconds
type Duration int64

const (
	Millisecond Duration = 1
	Second      Duration = 1000 * Millisecond
	Minute      Duration = 60 * Second
	Hour        Duration = 60 * Minute
)

func FromDuration(d time.Duration) Duration {
	return Duration(d / time.Millisecond)
}

func (d Duration) String() string {
	neg := ""
	if d < 0 {
		neg = "-"
		d = -d
	}

	hours := d / Hour
	d -= hours * Hour

	mins := d / Minute
	d -= mins * Minute

	secs := d / Second
	d -= secs * Second

	millis := d / Millisecond

	if hours > 0 {
		return fmt.Sprintf("%s%dh%dm%d.%03ds", neg, hours, mins, secs, millis)
	}
	if mins > 0 {
		return fmt.Sprintf("%s%dm%d.%03ds", neg, mins, secs, millis)
	}
	return fmt.Sprintf("%s%d.%03ds", neg, secs, millis)
}
