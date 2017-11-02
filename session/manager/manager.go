package manager

import (
	"ems/session"
	"ems/session/gorilla"

	"github.com/gorilla/sessions"
	"ems/middlewares"
	"net/http"
)
/*
usage:
	{SessionManager:  manager.SessionManager}
func (redirectBack *RedirectBack) Middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		returnTo := redirectBack.config.SessionManager.Get(req, "return_to")
		req = req.WithContext(context.WithValue(req.Context(), returnToKey, returnTo))

		if !redirectBack.Ignore(req) && returnTo != req.URL.String() {
			redirectBack.config.SessionManager.Add(w, req, "return_to", req.URL.String())
		}

		handler.ServeHTTP(w, req)
	})
}
 */
var SessionManager session.ManagerInterface = gorilla.New("_session", sessions.NewCookieStore([]byte("secret")))
func init() {
	middlewares.Use(middlewares.Middleware{
		Name: "session",
		Handler: func(handler http.Handler) http.Handler {
			return SessionManager.Middleware(handler)
		},
	})
}
