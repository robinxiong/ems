package routes

import (
	"ems/core/utils"
	"ems/site/app/controllers"
	"ems/site/db"
	"ems/wildcard_router"
	"net/http"

	"ems/core"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

var rootMux *http.ServeMux
var WildcardRouter *wildcard_router.WildcardRouter

func Router() *http.ServeMux {
	if rootMux == nil {
		router := chi.NewRouter()
		router.Use(middleware.Logger)
		router.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				var (
					tx          = db.DB
					siteContext = &core.Context{Request: req, Writer: w}
				)

				if locale := utils.GetLocale(siteContext); locale != "" {
					tx = tx.Set("l10n:locale", locale)
				}
				//todo publish2
				next.ServeHTTP(w, req)
			})

		})

		//routes
		router.Get("/", controllers.HomeIndex)

		rootMux = http.NewServeMux()
		WildcardRouter = wildcard_router.New()
		WildcardRouter.MountTo("/", rootMux)
		WildcardRouter.AddHandler(router)
	}
	return rootMux
}
