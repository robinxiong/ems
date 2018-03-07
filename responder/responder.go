package responder

import (
	"mime"
	"strings"
	"path/filepath"
	"net/http"
)

func Register(mimeType string, format string){
	//向mime中添加一个extension type, 默认从以下地方自动加载
	/*
	/etc/mime.types
	/etc/apache2/mime.types
	/etc/apache/mime.types
	 */
	mime.AddExtensionType("."+strings.TrimPrefix(format, "."), mimeType)
}

func init() {
	for mimeType, format := range map[string]string{
		"text/html":        "html",
		"application/json": "json",
		"application/xml":  "xml",
	} {
		Register(mimeType, format)
	}
}

type Responder struct {
	responds         map[string]func()
	DefaultResponder func()
}


// With could be used to register response handler for mime type formats, the formats could be string or []string
//     responder.With("html", func() {
//       writer.Write([]byte("this is a html request"))
//     }).With([]string{"json", "xml"}, func() {
//       writer.Write([]byte("this is a json or xml request"))
//     })
func With(formats interface{}, fc func()) *Responder {
	rep := &Responder{responds: map[string]func(){}}
	return rep.With(formats, fc)
}


// With could be used to register response handler for mime type formats, the formats could be string or []string
// 第一次注册的为默认的处理器
func (rep *Responder) With(formats interface{}, fc func()) *Responder {
	if f, ok := formats.(string); ok {
		rep.responds[f] = fc
	} else if fs, ok := formats.([]string); ok {
		for _, f := range fs {
			rep.responds[f] = fc
		}
	}

	if rep.DefaultResponder == nil {
		rep.DefaultResponder = fc
	}
	return rep
}


// 根据request访问对应的内容
func (rep *Responder) Respond(request *http.Request) {
	// get request format from url
	if ext := filepath.Ext(request.URL.Path); ext != "" {
		if respond, ok := rep.responds[strings.TrimPrefix(ext, ".")]; ok {
			respond()
			return
		}
	}

	// get request format from Accept
	for _, accept := range strings.Split(request.Header.Get("Accept"), ",") {
		if exts, err := mime.ExtensionsByType(accept); err == nil {
			for _, ext := range exts {
				if respond, ok := rep.responds[strings.TrimPrefix(ext, ".")]; ok {
					respond()
					return
				}
			}
		}
	}

	// 如果都没有找到，使用第一个With注册的response
	if rep.DefaultResponder != nil {
		rep.DefaultResponder()
	}
	return
}