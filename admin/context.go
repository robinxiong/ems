package admin

import (
	"net/http"

	"fmt"

	"html/template"

	"bytes"

	"ems/core"

	"github.com/qor/qor/utils"
)

//Context admin context, 用于 admin controller
type Context struct {
	//包含request, response, db, config, currentUser, role, errors
	*core.Context
	Admin        *Admin
	Settings     map[string]interface{}
	Resource     *Resource     //每一个请求包含对应的资源
	RouteHandler *routeHandler //route serveHTTP
	Action       string        //Controller中的action名字, 如果没有指定，则是context.Execute中的第一个参数比如dashboard
	Content      template.HTML //context.Render后的html
	Result       interface{}   //传递给template的数据
}

func (admin *Admin) NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{Context: &core.Context{Config: &core.Config{DB: admin.DB}, Request: r, Writer: w}, Admin: admin, Settings: map[string]interface{}{}}

}

func (context *Context) clone() *Context {
	return &Context{
		Context:  context.Context,
		Resource: context.Resource,
		Admin:    context.Admin,
		Result:   context.Result,
		Content:  context.Content,
		Settings: context.Settings,
		Action:   context.Action,
	}
}

/*
 Execute先读取layout.tmpl, 然后在调用context.Render渲染出具体的页面html, 然后在执行layout模板，并传递context作为数据对像
*/
func (context *Context) Execute(name string, result interface{}) {
	var tmpl *template.Template

	//todo: show to edit
	if context.Action == "" {
		context.Action = name
	}

	if content, err := context.Asset("layout.tmpl"); err == nil {
		if tmpl, err = template.New("layout").Funcs(context.FuncMap()).Parse(string(content)); err == nil {
			//解析layout.tpml中的header, footer
			for _, name := range []string{"header", "footer"} {
				if tmpl.Lookup(name) == nil {
					if content, err := context.Asset(name + ".tmpl"); err == nil {
						tmpl.Parse(string(content))
					}
				} else {
					utils.ExitWithMsg(err)
				}
			}
		} else {
			utils.ExitWithMsg(err)
		}

	}

	context.Result = result

	//渲染模板和数据，返回html
	context.Content = context.Render(name, result)
	if err := tmpl.Execute(context.Writer, context); err != nil {
		utils.ExitWithMsg(err)
	}

}

func (context *Context) Render(name string, results ...interface{}) template.HTML {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("Get error when render file %v: %v", name, r)
			utils.ExitWithMsg(err)
		}
	}()
	clone := context.clone()
	if len(results) > 0 {
		clone.Result = results[0]
	}
	return clone.renderWith(name, clone)
}

//获取指定的assets文件, 比如"layout.tmpl"
func (context *Context) Asset(layouts ...string) ([]byte, error) {
	var themes []string
	if context.Request != nil {
		if theme := context.Request.URL.Query().Get("theme"); theme != "" {
			themes = append(themes, theme)
		}
	}

	//todo: find the resource theme

	for _, layout := range layouts {
		//todo: find prefixes
		if content, err := context.Admin.AssetFS.Asset(layout); err == nil {
			return content, nil
		}
	}
	return []byte(""), fmt.Errorf("template not found: %v", layouts)
}
func (context *Context) resourcePath() string {
	if context.Resource == nil {
		return ""
	}
	return context.Resource.ToParam()
}

// renderWith render template based on data
func (context *Context) renderWith(name string, data interface{}) template.HTML {
	var (
		err     error
		content []byte
	)

	if content, err = context.Asset(name + ".tmpl"); err == nil {
		return context.renderText(string(content), data)
	}
	return template.HTML(err.Error())
}

// renderText render text based on data
func (context *Context) renderText(text string, data interface{}) template.HTML {
	var (
		err    error
		tmpl   *template.Template
		result = bytes.NewBufferString("")
	)

	if tmpl, err = template.New("").Funcs(context.FuncMap()).Parse(text); err == nil {
		if err = tmpl.Execute(result, data); err == nil {
			return template.HTML(result.String())
		}
	}

	return template.HTML(err.Error())
}
