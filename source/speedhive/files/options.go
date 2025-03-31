package files

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"zuercher.us/lapcharts/mode"
	"zuercher.us/lapcharts/source"
	"zuercher.us/lapcharts/util/flagx"
)

type Options struct {
	Files []string
}

var _ source.Options = &Options{}

func (o *Options) SourceMode() mode.SourceMode { return mode.SpeedhiveFiles }

func (o *Options) ConfigureFlags() {
	flag.Var(flagx.NewStringArrayValue(&o.Files, false), "files", "Speedhive Files: file names")
}

func (o *Options) Validate() error {
	if len(o.Files) < 2 {
		return errors.New("at least two files are required")
	}

	for _, file := range o.Files {
		if _, err := os.Stat(file); err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("file %q does not exist", file)
			}
			return fmt.Errorf("invalid file %q: %w", file, err)
		}
	}

	return nil
}

func (o *Options) Source() (source.Source, error) {
	return NewSource(o.Files...)
}
