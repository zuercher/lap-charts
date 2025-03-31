package mode

import "fmt"

type RenderMode string

const (
	Columns        RenderMode = "columns"
	HtmlLapsByTime RenderMode = "laps-by-time"
	HtmlDeltaByLap RenderMode = "delta-by-lap"
)

func AllRenderModes() []RenderMode {
	return []RenderMode{Columns, HtmlLapsByTime, HtmlDeltaByLap}
}

func ParseRenderMode(s string) (RenderMode, error) {
	for _, mode := range AllRenderModes() {
		if mode.String() == s || mode.ShortString() == s {
			return mode, nil
		}
	}

	return "", fmt.Errorf("unknown render mode %q", s)
}

func (m RenderMode) String() string { return string(m) }
func (m RenderMode) ShortString() string {
	switch m {
	case Columns:
		return "cols"
	case HtmlLapsByTime:
		return "laps"
	case HtmlDeltaByLap:
		return "delta"
	default:
		panic(fmt.Sprintf("unknown render mode %q", m))
	}
}
