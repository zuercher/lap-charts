package mode

import "fmt"

type SourceMode string

const (
	SpeedhiveFiles    SourceMode = "speedhive-files"
	SpeedhiveDownload SourceMode = "speedhive-download"
)

func AllSourceModes() []SourceMode {
	return []SourceMode{SpeedhiveFiles, SpeedhiveDownload}
}

func ParseSourceMode(s string) (SourceMode, error) {
	for _, mode := range AllSourceModes() {
		if mode.String() == s || mode.ShortString() == s {
			return mode, nil
		}
	}

	return "", fmt.Errorf("unknown mode %q", s)
}

func (m SourceMode) String() string { return string(m) }
func (m SourceMode) ShortString() string {
	switch m {
	case SpeedhiveFiles:
		return "shf"
	case SpeedhiveDownload:
		return "shd"
	default:
		panic(fmt.Sprintf("unknown mode %q", m))
	}
}
