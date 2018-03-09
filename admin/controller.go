package admin

import (
	"net/http"
	"strings"
	"time"
	"fmt"
	"crypto/md5"
	"mime"
	"path/filepath"
	"log"
)

//Controller包含admin对像，它在config中创建, 是一个全局的对像
type Controller struct {
	*Admin
	action *Action
}

func (ac *Controller) Dashboard(context *Context) {
	context.Execute("dashboard", nil)
}

//Show 渲染show page 单个
//Usage route.go RegisterResourceRoutes
func (ac *Controller) Show(context *Context) {
	log.Println("show")
}

//产品列表，多个值
func (ac *Controller) Index(context *Context) {
	log.Println("index")
}
var (
	cacheSince = time.Now().Format(http.TimeFormat)
)

func (ac *Controller) Asset(context *Context) {
	// /admin/assets/stylesheets/admin_default.css
	// /assets/stylesheets/admin_default.css
	file := strings.TrimPrefix(context.Request.URL.Path, ac.GetRouter().Prefix)
	//如果http请求中，带有If-Modified-Since，而它的值刚好等于服务器启动时间，则使用浏览器缓存的文件
	if context.Request.Header.Get("If-Modified-Since") == cacheSince {
		context.Writer.WriteHeader(http.StatusNotModified)
		return
	}
	//context.Writer.Header().Set("Last-Modified", cacheSince)

	//context在指定的路径下查找文件, 通常为
	// site/app/views，或者 ems/admin/views, 具体的路径设置，查看admin.SetAssetFS
	// admin/views/assets/stylesheets/admin_default.css
	if content, err := context.Asset(file); err == nil {
		etag := fmt.Sprintf("%x", md5.Sum(content))
		if context.Request.Header.Get("If-None-Match") == etag {
			context.Writer.WriteHeader(http.StatusNotModified)
			return
		}


		if ctype := mime.TypeByExtension(filepath.Ext(file)); ctype != "" {
			context.Writer.Header().Set("Content-Type", ctype)
		}

		//缓存300 second
		//context.Writer.Header().Set("Cache-control", "private, must-revalidate, max-age=300")
		//context.Writer.Header().Set("ETag", etag)
		context.Writer.Write(content)

	} else {

		http.NotFound(context.Writer, context.Request)

	}
}
