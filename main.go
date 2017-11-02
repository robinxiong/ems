package main

import (
	"net/http"
	"github.com/gorilla/sessions"

	"log"
	"os"
)

var store = sessions.NewFilesystemStore("", []byte("something-very-secret"))

func handler(w http.ResponseWriter, r *http.Request){
	log.Println(os.TempDir())
	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set some session values.
	session.Values["foo"] = "1"
	session.Values[42] = 43
	// Save it before we write to the response/return from the handler.
	session.Save(r, w)
}
func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
