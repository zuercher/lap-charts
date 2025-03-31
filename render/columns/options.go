package columns

import (
	"zuercher.us/lapcharts/mode"
	"zuercher.us/lapcharts/render"
	"zuercher.us/lapcharts/source"
)

type Options struct {
}

func (o *Options) RenderMode() mode.RenderMode {
	return mode.Columns
}
func (o *Options) ConfigureFlags() {
	// No flags to configure
}
func (o *Options) Validate() error {
	// No validation needed
	return nil
}
func (o *Options) Renderer(source source.Source) (render.Renderer, error) {
	// Create and return a new renderer based on the source
	return &columnRenderer{source}, nil
}
