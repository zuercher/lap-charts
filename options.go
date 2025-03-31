package lapcharts

import (
	"flag"
	"fmt"

	"zuercher.us/lapcharts/mode"
	"zuercher.us/lapcharts/render"
	"zuercher.us/lapcharts/render/columns"
	"zuercher.us/lapcharts/render/html"
	"zuercher.us/lapcharts/source"
	"zuercher.us/lapcharts/source/speedhive/download"
	"zuercher.us/lapcharts/source/speedhive/files"
	"zuercher.us/lapcharts/util/flagx"
)

type Options struct {
	SourceMode mode.SourceMode
	RenderMode mode.RenderMode

	output string

	shFileOptions     files.Options
	shDownloadOptions download.Options

	columnRenderOptions   columns.Options
	htmlLapsByTimeOptions html.LapsByTimeOptions
	htmlDeltaByLapOptions html.DeltaByLapOptions
}

func (o *Options) nestedSourceOptions() []source.Options {
	return []source.Options{
		&o.shFileOptions,
		&o.shDownloadOptions,
	}
}

func (o *Options) nestedRenderOptions() []render.Options {
	return []render.Options{
		&o.columnRenderOptions,
		&o.htmlLapsByTimeOptions,
		&o.htmlDeltaByLapOptions,
	}
}

func (o *Options) ConfigureFlags() {
	scv := &flagx.ChoiceValue[mode.SourceMode]{
		Value:   &o.SourceMode,
		Choices: mode.AllSourceModes(),
		Parse:   mode.ParseSourceMode,
	}
	flag.Var(
		scv,
		"source",
		"Choose the source mode. One of "+scv.Describe()+".",
	)
	for _, opts := range o.nestedSourceOptions() {
		opts.ConfigureFlags()
	}
	rcv := &flagx.ChoiceValue[mode.RenderMode]{
		Value:   &o.RenderMode,
		Choices: mode.AllRenderModes(),
		Parse:   mode.ParseRenderMode,
	}
	flag.Var(
		rcv,
		"render",
		"Choose the render mode. One of "+rcv.Describe()+".",
	)
	for _, opts := range o.nestedRenderOptions() {
		opts.ConfigureFlags()
	}

	flag.StringVar(&o.output, "output", "-", `Output path. Use "-" for stdout.`)
}

func (o *Options) Validate() error {
	for _, opts := range o.nestedSourceOptions() {
		if o.SourceMode == opts.SourceMode() {
			if err := opts.Validate(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (o *Options) Source() (source.Source, error) {
	for _, opts := range o.nestedSourceOptions() {
		if o.SourceMode == opts.SourceMode() {
			return opts.Source()
		}
	}
	return nil, fmt.Errorf("no source configured for mode %s", o.SourceMode)
}

func (o *Options) Renderer(src source.Source) (render.Renderer, error) {
	for _, opts := range o.nestedRenderOptions() {
		if o.RenderMode == opts.RenderMode() {
			return opts.Renderer(src)
		}
	}
	return nil, fmt.Errorf("no renderer configured for mode %s", o.RenderMode)
}
