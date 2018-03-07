package main

import (
	"ems/site/config"
	"flag"
	"fmt"
	"net/http"
	"os"
	"github.com/fatih/color"

	"ems/admin"

	"ems/publish2"
	"ems/site/app/static"
	"ems/site/config/application"

	"ems/site/db"
	"ems/site/utils/funcmapmaker"
	"path/filepath"
	adminapp "ems/site/app/admin"
	"github.com/go-chi/chi"
	"ems/core/utils"
	"github.com/go-chi/chi/middleware"
	"ems/site/config/auth"
	"ems/site/app/account"
	"ems/site/config/bindatafs"
)

func main() {

	// ./main --complie-templates true
	cmdLine := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	compileTemplate := cmdLine.Bool("compile-templates", false, "Compile Templates")
	cmdLine.Parse(os.Args[1:])

	var (
		Router = chi.NewRouter()
		Admin  = admin.New(&admin.AdminConfig{
			SiteName: "Site Name",  //影响site_name.css和site_name.js文件的加载
			Auth:  auth.AdminAuth{},
			DB:       db.DB.Set(publish2.VisibleMode, publish2.ModeOff).Set(publish2.ScheduleMode, publish2.ModeOff),
		})
		Application = application.New(&application.Config{
			Router: Router,
			Admin:  Admin,
			DB:     db.DB,
		})
	)

	funcmapmaker.AddFuncMapMaker(auth.Auth.Config.Render)

	Router.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			// for demo, don't use this for your production site
			w.Header().Add("Access-Control-Allow-Origin", "*")
			handler.ServeHTTP(w, req)
		})
	})


	Router.Use(middleware.Logger)
	Router.Use(middleware.Recoverer)


	Application.Use(account.New(&account.Config{}))
	Application.Use(adminapp.New(&adminapp.Config{}))
	Application.Use(static.New(&static.Config{
		Prefixs: []string{"/system"},
		Handler: utils.FileServer(http.Dir(filepath.Join(config.Root, "public"))),
	}))
	Application.Use(static.New(&static.Config{
		Prefixs: []string{"javascripts", "stylesheets", "images", "dist", "fonts", "vendors"},
		Handler: bindatafs.AssetFS.FileServer(http.Dir("public"), "javascripts", "stylesheets", "images", "dist", "fonts", "vendors"),
	}))

	if *compileTemplate {

	} else {
		fmt.Print(color.GreenString(fmt.Sprintf("Listening on: %v\n", config.Config.Port)))
		//config/config.go中配置了session middlewares
		if err := http.ListenAndServe(fmt.Sprintf(":%d", config.Config.Port), Application.NewServeMux()); err != nil {
			panic(err)
		}
		/*	if err := http.ListenAndServeTLS(fmt.Sprintf(":%d", config.Config.Port),"server.cert", "server.key", middlewares.Apply(mux)); err != nil {
			panic(err)
		}*/

	}

}
