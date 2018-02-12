package auth

import (
	"ems/auth/claims"
	"net/http"
	"path"
	"strings"
)

//serveMux用于将auth的路由添加到其它package, 比如site/config/routes/routes
type serveMux struct {
	*Auth
}

func (auth *Auth) NewServeMux() http.Handler {
	return &serveMux{auth}
}

func (serveMux *serveMux) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var (
		data    *claims.Claims
		reqPath = strings.TrimPrefix(req.URL.Path, serveMux.URLPrefix) //URLPrefix是在auth的字段，默认为auth
		paths   = strings.Split(reqPath, "/")
		context = &Context{Auth: serveMux.Auth, Claims: data, Request: req, Writer: w}
	)

	if len(paths) >= 2 {
		if paths[0] == "assets" {
			DefaultAssetHandler(context)
		}
	} else if len(paths) == 1 {
		switch paths[0] {
		case "login":
			serveMux.Auth.Render.Execute("auth/login", context, req, w) //context是传递给template的对像
		default:
			http.NotFound(w, req)
		}
	}

}

// AuthURL generate URL for auth
// 在auth_themes/clean/views/login.tmpl中，访问context下的.Auth(内嵌了auth)方法，获取/auth/assets/qor_auth.css
func (auth *Auth) AuthURL(pth string) string {
	return path.Join(auth.URLPrefix, pth)
}
