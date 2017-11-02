package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {

	// ./main --complie-templates true
	cmdLine := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	compileTemplate := cmdLine.Bool("compile-templates", false, "Compile Template")
	cmdLine.Parse(os.Args[1:])

	//log.Fatal(http.ListenAndServe(":8080", nil))
	log.Fatal(http.ListenAndServeTLS(":8080", "server.cert", "server.key", nil))
	if *compileTemplate {

	}
}
