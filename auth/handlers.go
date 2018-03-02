package auth

//controller中调用
import (
	"crypto/md5"
	"fmt"
	"mime"
	"net/http"
	"path"
	"path/filepath"
	"strings"
	"time"
	"ems/auth/claims"
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
func DefaultLoginhandler (context *Context, authorize func(*Context) (*claims.Claims, error)){
	var (
		req = context.Request
		w = context.Writer
		claims, err = authorize(context)  //调用provider的认证方法
	)
	//验证帐号密码成功
	if err == nil && claims != nil {

	}

	//todo:add session store
	//context.SessionStorer调用的是auth.SessionStorer
}