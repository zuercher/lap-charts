package render

import (
	"zuercher.us/lapcharts/mode"
	"zuercher.us/lapcharts/source"
)

type Options interface {
	RenderMode() mode.RenderMode
	ConfigureFlags()
	Validate() error
	Renderer(source.Source) (Renderer, error)
}
