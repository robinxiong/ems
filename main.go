package main

import (

	"log"
	"net/http"
)

func main() {
	var peoples *[]int = &[]int{1,2,3,}
	log.Println(peoples)
	http.HandleFunc("/", func( w http.ResponseWriter, req *http.Request){
		req.ParseForm()
		scopes := req.Form["scopes"]
		log.Println(scopes)
	})
	http.ListenAndServe(":8000", nil)

}
