package download

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"strconv"

	"zuercher.us/lapcharts/mode"
	"zuercher.us/lapcharts/source"
	"zuercher.us/lapcharts/source/speedhive/files"
	"zuercher.us/lapcharts/util/flagx"
)

type Options struct {
	EventURL  string
	Positions []int
}

var _ source.Options = &Options{}

func (o *Options) SourceMode() mode.SourceMode { return mode.SpeedhiveDownload }

func (o *Options) ConfigureFlags() {
	flag.StringVar(&o.EventURL, "event-url", "", "Speedhive Download: event URL")
	flag.Var(&flagx.ArrayValue[int]{Values: &o.Positions, Parse: strconv.Atoi}, "positions", "Speedhive Download: positions")
}

func (o *Options) Validate() error {
	if o.EventURL == "" {
		return errors.New("event-url is required")
	}

	u, err := url.Parse(o.EventURL)
	if err != nil {
		return fmt.Errorf("event-url is invalid: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("event-url scheme is not http or https: %s", u.Scheme)
	}

	if len(o.Positions) < 2 {
		return errors.New("at least two positions are required")
	}

	for _, pos := range o.Positions {
		if pos < 1 {
			return fmt.Errorf("position contains %d, must be 1 or higher", pos)
		}
	}

	return nil
}

func (o *Options) Source() (source.Source, error) {
	//	TODO: download files
	return files.NewSource()
}
