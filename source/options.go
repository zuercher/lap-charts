package source

import "zuercher.us/lapcharts/mode"

type Options interface {
	SourceMode() mode.SourceMode
	ConfigureFlags()
	Validate() error
	Source() (Source, error)
}
