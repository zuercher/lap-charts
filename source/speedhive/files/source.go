package files

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"zuercher.us/lapcharts/source"
	"zuercher.us/lapcharts/util/csvx"
	"zuercher.us/lapcharts/util/timex"
)

const (
	lapIdx      = 0
	posIdx      = 1
	timeIdx     = 2
	diffToP1Idx = 6
	speedIdx    = 7
)

func NewSource(files ...string) (source.Source, error) {
	s := &fileSource{
		num: len(files),
	}

	for _, file := range files {
		csv, err := csvx.New(file)
		if err != nil {
			return nil, fmt.Errorf("error loading %s: %w", file, err)
		}

		laps := source.NewLaps(file)

		firstRecord := true
		for record, err := range csv.Rows() {
			if err != nil {
				return nil, fmt.Errorf("error reading %s: %w", file, err)
			}
			if firstRecord {
				// First Row, read headers.
				if err := checkColumnIndices(record); err != nil {
					return nil, fmt.Errorf("error reading %s: %w", file, err)
				}
				firstRecord = false
				continue
			}

			lapStr := record[lapIdx]
			posStr := record[posIdx]
			var timeStr, speedStr string
			if strings.Contains(posStr, ":") {
				// Last row always has missing (not empty) pos column, so adjust the rest.
				timeStr = record[timeIdx-1]
				speedStr = record[speedIdx-1]
				posStr = ""
			} else {
				timeStr = record[timeIdx]
				speedStr = record[speedIdx]
			}

			lap, err := newLap(lapStr, posStr, timeStr, speedStr)
			if err != nil {
				return nil, fmt.Errorf("error parsing lap %d: %w", laps.NumLaps(), err)
			}

			if lap.Num == 1 {

			}

			last, ok := laps.LastLap()
			if ok && last.Num == lap.Num {
				// Ignore duplicate lap
				continue
			}

			if lap.Num == 1 {
				offsetStr := record[diffToP1Idx]
				offset, err := parseLapTime(offsetStr)
				if err == nil {
					laps.SetStartOffset(offset)
				}
			}
			laps.AppendLap(lap)
		}
		s.laps = append(s.laps, laps)

		fmt.Fprintf(os.Stderr, "Loaded %d laps from %s\n", laps.NumLaps(), file)
	}

	return s, nil
}

func checkColumnIndices(record []string) error {
	for idx, header := range record {
		switch strings.ToLower(header) {
		case "lap":
			if idx != lapIdx {
				return fmt.Errorf("unexpected lap column header at %d, expected %d", idx, lapIdx)
			}
		case "pos":
			if idx != posIdx {
				return fmt.Errorf("unexpected pos column header at %d, expected %d", idx, posIdx)
			}
		case "lap time":
			if idx != timeIdx {
				return fmt.Errorf("unexpected lap time column header at %d, expected %d", idx, timeIdx)
			}
		case "speed":
			if idx != speedIdx {
				return fmt.Errorf("unexpected speed column header at %d, expected %d", idx, speedIdx)
			}
		}
	}
	return nil
}

func newLap(lapStr, posStr, timeStr, speedStr string) (source.Lap, error) {
	lap, err := strconv.Atoi(lapStr)
	if err != nil {
		return source.Lap{}, fmt.Errorf("invalid lap number %q: %w", lapStr, err)
	}
	pos := -1
	if posStr != "" {
		pos, err = strconv.Atoi(posStr)
		if err != nil {
			return source.Lap{}, fmt.Errorf("invalid position %q: %w", posStr, err)
		}
	}

	lapTime, err := parseLapTime(timeStr)
	if err != nil {
		return source.Lap{}, err
	}

	speed, err := parseSpeed(speedStr)

	return source.Lap{
		Num:      lap,
		Position: pos,
		Time:     lapTime,
		SpeedKPH: speed,
	}, nil
}

func parseLapTime(timeStr string) (timex.Duration, error) {
	lapTime := timex.Duration(0)
	parts := strings.SplitN(timeStr, ".", 2)

	// h:m:s.000
	if len(parts) == 2 {
		digits := len(parts[1])
		if digits > 9 {
			return 0, fmt.Errorf("invalid lap time %q (fractional part): too many digits", timeStr)
		}

		frac, err := strconv.Atoi(parts[1])
		if err != nil {
			return 0, fmt.Errorf("invalid lap time %q (fractional part): %w", timeStr, err)
		}

		for digits%3 != 0 {
			frac *= 10
			digits++
		}
		switch digits {
		case 3:
			lapTime = timex.FromDuration(time.Duration(frac) * time.Millisecond)
		case 6:
			lapTime = timex.FromDuration(time.Duration(frac) * time.Microsecond)
		case 9:
			lapTime = timex.FromDuration(time.Duration(frac) * time.Nanosecond)
		default:
			return 0, fmt.Errorf("invalid lap time %q (fractional conversion): %d digits", timeStr, digits)
		}
	}

	// h:m:s
	parts = strings.Split(parts[0], ":")
	if len(parts) > 3 {
		return 0, fmt.Errorf("invalid lap time %q (too many parts): %d", timeStr, len(parts))
	}
	mult := 1
	for len(parts) > 0 {
		n, err := strconv.Atoi(parts[len(parts)-1])
		if err != nil {
			return 0, fmt.Errorf("invalid lap time %q: %w", timeStr, err)
		}
		lapTime += timex.Duration(n*mult) * timex.Second
		mult *= 60
		parts = parts[:len(parts)-1]
	}
	return lapTime, nil
}

func parseSpeed(speedStr string) (float64, error) {
	conversions := []struct {
		suffix string
		factor float64
	}{
		{"km/h", 1},
		{"kph", 1},
		{"m/h", 1.609344},
		{"mph", 1.609344},
	}

	for _, conv := range conversions {
		if s := strings.TrimSuffix(speedStr, conv.suffix); s != speedStr {
			speed, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
			if err != nil {
				return 0, fmt.Errorf("invalid speed %q: %w", speedStr, err)
			}
			return speed * conv.factor, nil
		}
	}
	// unknown conversion, silently ignore the speed
	return 0, nil
}

type fileSource struct {
	num  int
	laps []source.Laps
}

func (s *fileSource) NumCars() int             { return s.num }
func (s *fileSource) Laps(car int) source.Laps { return s.laps[car] }
