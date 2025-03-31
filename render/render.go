package render

import "zuercher.us/lapcharts/util/iox"

type Renderer interface {
	Render(iox.Writer) error
}
