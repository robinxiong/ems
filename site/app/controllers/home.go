package controllers

import (
	"ems/site/config"
	"net/http"
)

func HomeIndex(w http.ResponseWriter, req *http.Request) {

	config.View.Execute("home_index", map[string]interface{}{}, req, w)
}
