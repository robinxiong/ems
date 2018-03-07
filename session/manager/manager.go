package manager

import (
	"ems/session"
	"ems/session/gorilla"

	"ems/middlewares"
	"net/http"

	"github.com/gorilla/sessions"
	"time"
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
//secureCoookie用来加密session的密码
var (
	key                                     = []byte("加密session为cookie的密码")
	cookieStore                             = sessions.NewCookieStore(key)
	SessionManager session.ManagerInterface = gorilla.New("_session", cookieStore)
)

func init() {
	cookieStore.Options = &sessions.Options{
		Path:   "/",
		MaxAge: int(time.Hour * 24), //0则表示只在当前有效
	}
	middlewares.Use(middlewares.Middleware{
		Name: "session",
		Handler: func(handler http.Handler) http.Handler {
			return SessionManager.Middleware(handler)
		},
	})
}
