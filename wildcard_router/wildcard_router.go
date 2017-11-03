package wildcard_router

import (
	"net/http"
	"strings"
)

//有的路由是根据数据库保存的url来决定的，所以没有一定的规则来定义url，所以wildcard_router是用来动态设置路由
//根据url查找每一个注册的handler中的数据库表，如果在某一个表中没有找到，则查找下一个数据库表，则到找到相应的路径
/*

	if !db.First(&page, "url = ?", req.URL.Path).RecordNotFound() {
		w.Write([]byte(page.Body))
	}
	if !db.First(&faq, "url = ?", req.URL.Path).RecordNotFound() {
		w.Write([]byte(fmt.Sprintf("%v: %v", faq.Question, faq.Answer)))
	}
 */
type WildcardRouter struct {
	middlewares []func(writer http.ResponseWriter, request *http.Request)
	handlers []http.Handler
	notFoundHandler http.HandlerFunc
}


func New() *WildcardRouter{
	return &WildcardRouter{}
}


func (w *WildcardRouter) MountTo(mountTo string, mux *http.ServeMux){
	mountTo = "/" + strings.Trim(mountTo, "/")
	mux.Handle(mountTo, w)
	mux.Handle(mountTo + "/", w)
}

func (w *WildcardRouter) AddHandler(handler http.Handler) {
	w.handlers = append(w.handlers, handler)
}

// NoRoute will set handler to handle 404
func (w *WildcardRouter) NoRoute(handler http.HandlerFunc) {
	w.notFoundHandler = handler
}

// Use will append new middleware, 仅在wildcardRouter中使用的中间件
func (w *WildcardRouter) Use(middleware func(writer http.ResponseWriter, request *http.Request)) {
	w.middlewares = append(w.middlewares, middleware)
}

func (w *WildcardRouter) ServeHTTP(writer http.ResponseWriter, req *http.Request){
	//它用于查找所有注册到wildcardRouter中的handler, 所以在没有找到的情况下，则跳过状态的设置
	wildcardRouterWriter := &WildcardRouterWriter{writer, 0, false}

	for _, middleware := range w.middlewares {
		middleware(writer, req)
	}


	for _, handler := range w.handlers {
		if handler.ServeHTTP(wildcardRouterWriter, req); wildcardRouterWriter.isProcessed() {
			return
		}
		wildcardRouterWriter.reset()
	}

	//如果没有找到，这时需要设置为true, 这样到跟将404状态码返回给客户端
	wildcardRouterWriter.skipNotFoundCheck = true
	if w.notFoundHandler != nil {
		w.notFoundHandler(writer, req)
	} else {
		http.NotFound(wildcardRouterWriter, req)
	}
}

//WildcardRouterWriter 用于捕获状态, 它实现了http.ResponseWriter
type WildcardRouterWriter struct {
	http.ResponseWriter
	//保存状态码
	status int
	//是否跳过状态检查, 如果指定为true, 则会在所有的状态码返回给客户端，如果为false, 则跳过404状态码的返回
	skipNotFoundCheck bool
}

func (w WildcardRouterWriter) Status() int {
	return w.status
}


// WriteHeader only set status code when not 404
func (w *WildcardRouterWriter) WriteHeader(statusCode int) {
	if w.skipNotFoundCheck || statusCode != http.StatusNotFound {
		w.ResponseWriter.WriteHeader(statusCode)
	}
	w.status = statusCode
}

// Write only set content when not 404
func (w *WildcardRouterWriter) Write(data []byte) (int, error) {
	if w.skipNotFoundCheck || w.status != http.StatusNotFound {
		w.status = http.StatusOK
		return w.ResponseWriter.Write(data)
	}
	return 0, nil
}

func (w *WildcardRouterWriter) reset() {
	w.skipNotFoundCheck = false
	w.Header().Set("Content-Type", "")
	w.status = 0
}

func (w WildcardRouterWriter) isProcessed() bool {
	return w.status != http.StatusNotFound && w.status != 0
}
