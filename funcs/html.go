package funcs

import (
	"html/template"
)

var pathCache = make(map[string]string)

func safeHTML(s string) template.HTML {
	return template.HTML(s)
}

func safeAttr(s string) template.HTMLAttr {
	return template.HTMLAttr(s)
}

func safeCSS(s string) template.CSS {
	return template.CSS(s)
}

func safeJS(s string) template.JS {
	return template.JS(s)
}

func safeURL(s string) template.URL {
	return template.URL(s)
}
