package admin

import (
	"ems/core"
	"ems/core/utils"
	"ems/roles"
	"log"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"sort"
	"strings"
	"time"
)

//admin是相关的路由

func (admin *Admin) NewServeMux(prefix string) http.Handler {
	//注册默认的路由和中间件, 路由在admin中指定，所以在site/config/admin中，会指定路由信息
	router := admin.router
	if router != nil {
		router.Prefix = prefix
	}
	adminController := &Controller{Admin: admin}
	router.Get("", adminController.Dashboard)
	router.Use(&Middleware{
		Name: "admin_handler",
		Handler: func(context *Context, middleware *Middleware) {
			context.Writer.Header().Set("Cache-control", "no-store")
			context.Writer.Header().Set("Pragma", "no-cache")
			if context.RouteHandler != nil {
				context.RouteHandler.Handle(context)
				return
			}
			http.NotFound(context.Writer, context.Request)
		},
	})

	return &serveMux{admin: admin}
}

func (admin *Admin) MountTo(mountTo string, mux *http.ServeMux) {
	prefix := "/" + strings.Trim(mountTo, "/")

	serveMux := admin.NewServeMux(prefix)

	mux.Handle(prefix, serveMux)
	mux.Handle(prefix+"/", serveMux)
}

type serveMux struct {
	admin *Admin
}

//ServerHTTP 匹配route中注册的路由，调用它的处理函数
func (serverMux *serveMux) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var (
		admin        = serverMux.admin
		RelativePath = "/" + strings.Trim(strings.TrimPrefix(req.URL.Path, admin.router.Prefix), "/")
		context      = admin.NewContext(w, req)
	)

	//parse request form, 2M
	req.ParseMultipartForm(2 * 1024 * 1024)
	if method := req.Form.Get("_method"); method != "" {
		req.Method = strings.ToUpper(method)
	}

	//如果请求的路径是以assets开头，则调用Controller中的Asset方法
	// /admin/assets/stylesheets/admin_default.css
	//RelativePath删除了prefix, /admin
	if regexp.MustCompile("^/assets/.*$").MatchString(RelativePath) && strings.ToUpper(req.Method) == "GET" {
		//调用controller中的Asset Action, 设置静态文件的缓存, 并向浏览器返回文件

		(&Controller{Admin: admin}).Asset(context)
		return
	}
	defer func() func() {
		begin := time.Now()
		return func() {
			log.Printf("Finish [%s] %s Took %.2fms\n", req.Method, req.RequestURI, time.Now().Sub(begin).Seconds()*1000)
		}
	}()()

	//设置当前用户，如果没有当前用户，则转向到登陆页面
	var currentUser core.CurrentUser //core/context.go
	//var permissionModel roles.Permission

	//它在config/admin/admin.go的init方法中初始化Auth: auth.AdminAuth{}
	//如果不需要认证，则在创建admin时，可以不指定Auth
	if admin.Auth != nil {
		if currentUser = admin.Auth.GetCurrentUser(context); currentUser == nil {
			http.Redirect(w, req, admin.Auth.LoginURL(context), http.StatusSeeOther)
			return
		}

		context.CurrentUser = currentUser
		context.SetDB(context.GetDB().Set("qor:current_user", context.CurrentUser))
	}

	context.Roles = roles.MatchedRoles(req, currentUser) //core/resource.go中定义了Roles, MatchedRoles是匹配req和当前用户，是否符合系统中定义的角色，符合返回角色名称

	//router中注册的controller
	handlers := admin.router.routers[strings.ToUpper(req.Method)]
	for _, handler := range handlers {
		//todo: hander的权限
		if params, _, ok := utils.ParamsMatch(handler.Path, RelativePath); ok {
			if len(params) > 0 {
				req.URL.RawQuery = url.Values(params).Encode() + "&" + req.URL.RawQuery
			}

			context.RouteHandler = handler
			context.setResource(handler.Config.Resource)
			break
		}
	}

	// Call first middleware
	// 调用 NewServeMux中的qor_handler, 它执行上一步获取到的handler, 即controller中的handle
	for _, middleware := range admin.router.middlewares {
		middleware.Handler(context, middleware)
		break
	}
}

//在admin.AddResource中调用， 将一个resource注册到路由
func (admin *Admin) RegisterResourceRouters(res *Resource, actions ...string) {
	var (
		primaryKeyParams = res.ParamIDName()         // :资源名_id
		adminController  = &Controller{Admin: admin} //为每个资源创建一个controller, 比如NewServeMux也创建了一个controller用于注册""路径
	)

	//遍历actions
	for _, action := range actions {
		switch strings.ToLower(action) {
		case "read":
			if res.Config.Singleton {
				//如果资源添加时为Singleton，单个
				res.RegisterRoute("GET", "/", adminController.Show, &RouteConfig{PermissionMode: roles.Read})
			} else {

				// Index, 多个
				res.RegisterRoute("GET", "/", adminController.Index, &RouteConfig{PermissionMode: roles.Read})

				// Show  显示单个 /:product_id/
				res.RegisterRoute("GET", primaryKeyParams, adminController.Show, &RouteConfig{PermissionMode: roles.Read})
			}
		}
	}
}

func (res *Resource) RegisterRoute(method string, relativePath string, handler requestHandler, config *RouteConfig) {
	if config == nil {
		config = &RouteConfig{}
	}

	config.Resource = res

	var (
		prefix string
		param  = res.ToParam()
		router = res.GetAdmin().router
	)

	//一个匿名函数
	if prefix = func(r *Resource) string {
		currentParam := param //当前resource的资源名称 colors or order_items
		for r.ParentResource != nil {
			parentPath := r.ParentResource.ToParam()
			//不向router中注册相同的resource
			if parentPath == param {
				return ""
			}
			currentParam = path.Join(parentPath, r.ParentResource.ParamIDName(), currentParam)
			r = r.ParentResource
		}
		return "/" + strings.Trim(currentParam, "/")
	}(res); prefix == "" {
		return
	}
	//relativePath 通常 "/", 即当前资源下的哪一个操作
	//Get, Post方法中会调用newHandle将，handler和config绑定，在ServeHTTP，执行路由时，通过找到handler就能获取它所对应的Resource
	switch strings.ToUpper(method) {
	case "GET":
		router.Get(path.Join(prefix, relativePath), handler, config)
	case "POST":
		router.Post(path.Join(prefix, relativePath), handler, config)
	case "PUT":
		router.Put(path.Join(prefix, relativePath), handler, config)
	case "DELETE":
		router.Delete(path.Join(prefix, relativePath), handler, config)
	}
}

//Middleware
type Middleware struct {
	Name    string
	Handler func(*Context, *Middleware)
	next    *Middleware
}

// Next will call the next middleware
func (middleware Middleware) Next(context *Context) {
	if next := middleware.next; next != nil {
		next.Handler(context, next)
	}
}

//Router
type Router struct {
	Prefix      string
	routers     map[string][]*routeHandler
	middlewares []*Middleware
}

func newRouter() *Router {
	return &Router{routers: map[string][]*routeHandler{
		"GET":    {},
		"PUT":    {},
		"POST":   {},
		"DELETE": {},
	}}
}

// Use reigster a middleware to the router
func (r *Router) Use(middleware *Middleware) {
	// compile middleware
	for index, m := range r.middlewares {
		// replace middleware have same name
		if m.Name == middleware.Name {
			middleware.next = m.next
			r.middlewares[index] = middleware
			if index > 1 {
				r.middlewares[index-1].next = middleware
			}
			return
		} else if len(r.middlewares) > index+1 {
			m.next = r.middlewares[index+1]
		} else if len(r.middlewares) == index+1 {
			m.next = middleware
		}
	}

	r.middlewares = append(r.middlewares, middleware)
}

// GetMiddleware get registered middleware
func (r *Router) GetMiddleware(name string) *Middleware {
	for _, middleware := range r.middlewares {
		if middleware.Name == name {
			return middleware
		}
	}
	return nil
}

// Get register a GET request handle with the given path
func (r *Router) Get(path string, handle requestHandler, config ...*RouteConfig) {
	r.routers["GET"] = append(r.routers["GET"], newRouteHandler(path, handle, config...))
	r.sortRoutes(r.routers["GET"])
}

// Post register a POST request handle with the given path
func (r *Router) Post(path string, handle requestHandler, config ...*RouteConfig) {
	r.routers["POST"] = append(r.routers["POST"], newRouteHandler(path, handle, config...))
	r.sortRoutes(r.routers["POST"])
}

// Put register a PUT request handle with the given path
func (r *Router) Put(path string, handle requestHandler, config ...*RouteConfig) {
	r.routers["PUT"] = append(r.routers["PUT"], newRouteHandler(path, handle, config...))
	r.sortRoutes(r.routers["PUT"])
}

// Delete register a DELETE request handle with the given path
func (r *Router) Delete(path string, handle requestHandler, config ...*RouteConfig) {
	r.routers["DELETE"] = append(r.routers["DELETE"], newRouteHandler(path, handle, config...))
	r.sortRoutes(r.routers["DELETE"])
}

var wildcardRouter = regexp.MustCompile(`/:\w+`)

func (r *Router) sortRoutes(routes []*routeHandler) {
	sort.SliceStable(routes, func(i, j int) bool {
		iIsWildcard := wildcardRouter.MatchString(routes[i].Path)
		jIsWildcard := wildcardRouter.MatchString(routes[j].Path)
		// i regexp (true), j static (false) => false
		// i static (true), j regexp (true) => true
		if iIsWildcard != jIsWildcard {
			return jIsWildcard
		}
		return len(routes[i].Path) > len(routes[j].Path)
	})
}
