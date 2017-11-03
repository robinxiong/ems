package main

import (
	"ems/site/config"
	"ems/site/config/routes"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/fatih/color"
	"ems/middlewares"
)

func main() {

	// ./main --complie-templates true
	cmdLine := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	compileTemplate := cmdLine.Bool("compile-templates", false, "Compile Template")
	cmdLine.Parse(os.Args[1:])

	mux := http.NewServeMux()
	mux.Handle("/", routes.Router())



	if *compileTemplate {

	}

	fmt.Print(color.GreenString(fmt.Sprintf("Listening on: %v\n", config.Config.Port)))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.Config.Port), middlewares.Apply(mux)); err != nil {
		panic(err)
	}
/*	if err := http.ListenAndServeTLS(fmt.Sprintf(":%d", config.Config.Port),"server.cert", "server.key", middlewares.Apply(mux)); err != nil {
		panic(err)
	}*/


}
