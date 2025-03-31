package lapcharts

import (
	"os"

	"zuercher.us/lapcharts/util/iox"
)

func Generate(o *Options) error {
	source, err := o.Source()
	if err != nil {
		return err
	}

	renderer, err := o.Renderer(source)
	if err != nil {
		return err
	}

	var output iox.Writer
	if o.output == "-" {
		output = iox.NewWriter(os.Stdout)
	} else {
		file, err := os.Create(o.output)
		if err != nil {
			return err
		}
		defer file.Close()
		output = iox.NewWriter(file)
	}

	return renderer.Render(output)
}
