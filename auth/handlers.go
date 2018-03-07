package auth

//controller中调用
import (
	"crypto/md5"
	"ems/auth/claims"
	"ems/responder"
	"ems/session"
	"fmt"
	"html/template"
	"mime"
	"net/http"
	"path"
	"path/filepath"
	"strings"
	"time"
)

var cacheSince = time.Now().Format(http.TimeFormat)

func DefaultAssetHandler(context *Context) {
	asset := strings.TrimPrefix(context.Request.URL.Path, context.Auth.URLPrefix) //去掉路径中的/auth
	if context.Request.Header.Get("If-Modified-Since") == cacheSince {
		context.Writer.WriteHeader(http.StatusNotModified)
		return
	}
	context.Writer.Header().Set("Last-Modified", cacheSince)

	if content, err := context.Render.Asset(path.Join("/auth", asset)); err == nil {
		etag := fmt.Sprintf("%x", md5.Sum(content))
		if context.Request.Header.Get("If-None-Match") == etag {
			context.Writer.WriteHeader(http.StatusNotModified)
			return
		}

		if ctype := mime.TypeByExtension(filepath.Ext(asset)); ctype != "" {

			context.Writer.Header().Set("Content-Type", ctype)
		}

		context.Writer.Header().Set("Cache-control", "private, must-revalidate, max-age=300")
		context.Writer.Header().Set("ETag", etag)
		context.Writer.Write(content)
	} else {

		http.NotFound(context.Writer, context.Request)
	}
}

//验证用户帐号和密码时调用,它在auth.New中赋值给auth.LoginHandler， 而auth.LoginHandler在provider的Login函数中调用
func DefaultLoginHandler(context *Context, authorize func(*Context) (*claims.Claims, error)) {
	var (
		req         = context.Request
		w           = context.Writer
		claims, err = authorize(context) //调用provider的认证方法
	)
	//验证帐号密码成功
	if err == nil && claims != nil {
		context.SessionStorer.Flash(w, req, session.Message{Message: "logged"})
		respondAfterLogged(claims, context)
		return
	}

	//context.SessionStorer调用的是auth.SessionStorer
	//写入flash到session, 可能会报错
	context.SessionStorer.Flash(w, req, session.Message{Message: template.HTML(err.Error()), Type: "error"})

	//向浏览器返回信息
	responder.With("html", func() {

		context.Auth.Config.Render.Execute("auth/login", context, req, w)
	}).With([]string{"json"}, func() {
		// TODO write json error
	}).Respond(context.Request)
}


//注册用户信息，它在password.Register中调用，第一个参数为context, 第二个参数为provider.RegisterHandler(负责保存用户到数据库, 以及对密码的加密等逻辑)
func DefaultRegisterHandler(context *Context, register func(*Context)(*claims.Claims, error)) {
	var (
		req         = context.Request
		w           = context.Writer
		claims, err = register(context)
	)

	if err == nil && claims != nil {
		respondAfterLogged(claims, context)
		return
	}
	//输出错误信息到session, 比如用户名已经存在
	context.SessionStorer.Flash(w, req, session.Message{Message: template.HTML(err.Error()), Type: "error"})

	// error handling
	responder.With("html", func() {
		context.Auth.Config.Render.Execute("auth/register", context, req, w)
	}).With([]string{"json"}, func() {
		// TODO write json error
	}).Respond(context.Request)
}



func respondAfterLogged(claims *claims.Claims, context *Context) {
	// 将claims写入到session
	context.Auth.Login(context.Writer, context.Request, claims)

	responder.With("html", func() {
		// write cookie, 转回到登录页面, 验证claims
		context.Auth.Redirector.Redirect(context.Writer, context.Request, "login")
	}).With([]string{"json"}, func() {
		// TODO write json token
	}).Respond(context.Request)
}