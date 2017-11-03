package controllers

import (
	"net/http"
	"fmt"
)

func HomeIndex(w http.ResponseWriter, req *http.Request){
	fmt.Fprintf(w, "home page")
}