package render

import (
	"html/template"
	"net/http"
)

type Template struct {
	render *Render
	layout string
	usingDefaultLayout bool
	funcMap template.FuncMap
}

//FuncMap  get func maps from tmpl
func (tmpl *Template) funcMapMaker(req *http.Request, writer http.ResponseWriter) template.FuncMap{
	var funcMap = template.FuncMap{}

	for key, fc := range tmpl.render.funcMaps {
		funcMap[key] = fc
	}

	if tmpl.render.Config.FuncMapMaker != nil {
		for key, fc := range tmpl.render.Config.FuncMapMaker(tmpl.render, req, writer) {
			funcMap[key] = fc
		}
	}

	for key, fc := range tmpl.funcMap {
		funcMap[key] = fc
	}

	return funcMap
}

// Funcs register Funcs for tmpl
func (tmpl *Template) Funcs(funcMap template.FuncMap) *Template {
	tmpl.funcMap = funcMap
	return tmpl
}


// Render render tmpl
func (tmpl *Template) Render(templateName string, obj interface{}, request *http.Request, writer http.ResponseWriter) (template.HTML, error) {

}