package html

import "html/template"

func templateFuncs() template.FuncMap {
	return map[string]any{
		"add": func(a, b int) int {
			return a + b
		},
	}
}
