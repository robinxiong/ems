package main

import (
	"ems/site/config"
	"ems/site/config/routes"
	"flag"
	"fmt"
	"net/http"
	"os"

	"ems/middlewares"

	"github.com/fatih/color"
	"ems/site/config/admin"
	"ems/render"
	"html/template"
	"ems/i18n/inline_edit"

	"ems/site/config/utils"
	"ems/site/config/i18n"
)

func main() {

	// ./main --complie-templates true
	cmdLine := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	compileTemplate := cmdLine.Bool("compile-templates", false, "Compile Template")
	cmdLine.Parse(os.Args[1:])

	mux := http.NewServeMux()
	mux.Handle("/", routes.Router())

	admin.Admin.MountTo("/admin", mux)


	/* 配置 funcMapMaker 它将在render的template执行时调用，然后返回模板的funcMap */
	config.View.FuncMapMaker = func(render *render.Render, req *http.Request, w http.ResponseWriter) template.FuncMap {
		funcMap := template.FuncMap{}
		//添加模板的中的t方法，用于语言的翻译
		for key, fc := range inline_edit.FuncMap(i18n.I18n, utils.GetCurrentLocale(req), false) {
			funcMap[key] = fc
		}
		return funcMap
	}


	if *compileTemplate {

	}

	fmt.Print(color.GreenString(fmt.Sprintf("Listening on: %v\n", config.Config.Port)))
	//config/config.go中配置了session middlewares
	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.Config.Port), middlewares.Apply(mux)); err != nil {
		panic(err)
	}
	/*	if err := http.ListenAndServeTLS(fmt.Sprintf(":%d", config.Config.Port),"server.cert", "server.key", middlewares.Apply(mux)); err != nil {
		panic(err)
	}*/

}
