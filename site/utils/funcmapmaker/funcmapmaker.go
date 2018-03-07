package funcmapmaker

import (
	"html/template"
	"net/http"

	"ems/render"
	"ems/i18n/inline_edit"


	"ems/site/utils"
	"ems/site/config/i18n"
	"ems/session"
	"ems/session/manager"
)

// GetEditMode get edit mode
func GetEditMode(w http.ResponseWriter, req *http.Request) bool {
	return false
}

// AddFuncMapMaker add FuncMapMaker to view
func AddFuncMapMaker(view *render.Render) *render.Render {
	oldFuncMapMaker := view.FuncMapMaker
	view.FuncMapMaker = func(render *render.Render, req *http.Request, w http.ResponseWriter) template.FuncMap {
		funcMap := template.FuncMap{}
		if oldFuncMapMaker != nil {
			funcMap = oldFuncMapMaker(render, req, w)
		}

		// Add `t` method
		for key, fc := range inline_edit.FuncMap(i18n.I18n, utils.GetCurrentLocale(req), GetEditMode(w, req)) {
			funcMap[key] = fc
		}


		funcMap["flashes"] = func() []session.Message {
			return manager.SessionManager.Flashes(w, req)
		}

		return funcMap
	}

	return view
}
