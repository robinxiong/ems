package admin

import (
	"bytes"
	"ems/core/utils"
	"ems/roles"
	"fmt"
	"html/template"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
)

//for context

func (context *Context) FuncMap() template.FuncMap {
	funcMap := template.FuncMap{
		"page_title":             context.pageTitle,
		"t":                      context.t,                    //翻译成本地语言，它调用admin.T
		"stylesheet_tag":         context.styleSheetTag,        //根据指定的字符串，返回template.HTML, 比如admin_default, 它返回/admin/assets/stylesheets/admin_default.css, 对于assets的请求，可以查看admin.Route中的serveHTTP
		"load_admin_stylesheets": context.loadAdminStyleSheets, //返回指定siteName的css, 如果没有指定siteName, 则返回application_css
		"load_theme_stylesheets": context.loadThemeStyleSheets,
		"javascript_tag":         context.javaScriptTag,  //从/admin/views/assets/javascripts中加载指定名字的js文件
		"load_actions":           context.loadActions, //加载views/action下的子模板, action会按照文件前面的数字进行排名
		"qor_theme_class":        context.themesClass, //返回theme style的类名称，用于body标签
		"render":                 context.Render,      //读取指定的模板
		"logout_url":             context.logoutURL,   //sidebar.tmpl获取logout url
		"get_menus":              context.getMenus,    //获取系统菜单，并且传递给sidebar.tmpl
		"link_to":                context.linkTo,      //翻译link_to的名称
		"load_admin_javascripts": context.loadAdminJavaScripts,  //根据site来加载js文件

	}
	return funcMap
}
func (context *Context) URLFor(value interface{}, resources ...*Resource) string {
	return value.(string)
}
func (context *Context) pageTitle() template.HTML {
	if context.Resource == nil {
		return context.t("qor_admin.layout.title", "Admin")
	}
	return ""
}

func (context *Context) styleSheetTag(names ...string) template.HTML {
	var results []string
	for _, name := range names {
		name = path.Join(context.Admin.GetRouter().Prefix, "assets", "stylesheets", name+".css")
		results = append(results, fmt.Sprintf(`<link type="text/css" rel="stylesheet" href="%s">`, name))
	}
	return template.HTML(strings.Join(results, ""))
}

//layout.tmpl加使用, 根据站点名称加载css, 这里为/admin/assets/stylesheets/qor_demo.css
func (context *Context) loadAdminStyleSheets() template.HTML {
	var siteName = context.Admin.SiteName
	if siteName == "" {
		siteName = "application"
	}
	var file = path.Join("assets", "stylesheets", strings.ToLower(strings.Replace(siteName, " ", "_", -1))+".css")
	if _, err := context.Asset(file); err == nil {
		return template.HTML(fmt.Sprintf(`<link type="text/css" rel="stylesheet" href="%s">`, path.Join(context.Admin.GetRouter().Prefix, file)))
	}
	return ""

}

func (context *Context) loadAdminJavaScripts() template.HTML {
	var siteName = context.Admin.SiteName
	if siteName == "" {
		siteName = "application"
	}

	var file = path.Join("assets", "javascripts", strings.ToLower(strings.Replace(siteName, " ", "_", -1))+".js")
	if _, err := context.Asset(file); err == nil {
		return template.HTML(fmt.Sprintf(`<script src="%s"></script>`, path.Join(context.Admin.GetRouter().Prefix, file)))
	}
	return ""
}


func (context *Context) loadThemeStyleSheets() template.HTML {
	var results []string
	for _, themeName := range context.getThemeNames() {
		var file = path.Join("themes", themeName, "assets", "stylesheets", themeName+".css")
		if _, err := context.Asset(file); err == nil {
			results = append(results, fmt.Sprintf(`<link type="text/css" rel="stylesheet" href="%s?theme=%s">`, path.Join(context.Admin.GetRouter().Prefix, "assets", "stylesheets", themeName+".css"), themeName))
		}
	}

	return template.HTML(strings.Join(results, " "))
}

func (context *Context) javaScriptTag(names ...string) template.HTML {
	var results []string
	for _, name := range names {
		name = path.Join(context.Admin.GetRouter().Prefix, "assets", "javascripts", name+".js")
		results = append(results, fmt.Sprintf(`<script src="%s"></script>`, name))
	}
	return template.HTML(strings.Join(results, ""))
}

func (context *Context) getThemeNames() (themes []string) {
	themesMap := map[string]bool{}

	if context.Resource != nil {
		for _, theme := range context.Resource.Config.Themes {
			if _, ok := themesMap[theme.GetName()]; !ok {
				themes = append(themes, theme.GetName())
			}
		}
	}

	return
}

// {{$actions := load_actions "header" }}
func (context *Context) loadActions(action string) template.HTML {
	var (
		actionPatterns, actionKeys, actionFiles []string
		actions                                 = map[string]string{}
	)
	switch action {
	case "index", "show", "edit", "new":
		//todo loadAction for index
	case "global":
		actionPatterns = []string{"actions/*.tmpl"}
	default:
		actionPatterns = []string{filepath.Join("actions", action, "*.tmpl")}
	}
	for _, pattern := range actionPatterns {
		for _, themeName := range context.getThemeNames() {
			if resourcePath := context.resourcePath(); resourcePath != "" {
				if matches, err := context.Admin.AssetFS.Glob(filepath.Join("themes", themeName, resourcePath, pattern)); err == nil {
					actionFiles = append(actionFiles, matches...)
				}
			}

			if matches, err := context.Admin.AssetFS.Glob(filepath.Join("themes", themeName, pattern)); err == nil {
				actionFiles = append(actionFiles, matches...)
			}
		}

		if resourcePath := context.resourcePath(); resourcePath != "" {
			if matches, err := context.Admin.AssetFS.Glob(filepath.Join(resourcePath, pattern)); err == nil {
				actionFiles = append(actionFiles, matches...)
			}
		}

		//找到header文件夹下的所有tmpl, 1.page_title.tmpl 5.userinfo.tmpl, 6.searchbar.tmpl
		if matches, err := context.Admin.AssetFS.Glob(pattern); err == nil {
			actionFiles = append(actionFiles, matches...)
		}
	}

	// before files have higher priority
	for _, actionFile := range actionFiles {
		//替换文件中的数字和点为空 1.
		base := regexp.MustCompile("^\\d+\\.").ReplaceAllString(path.Base(actionFile), "")
		//如果actions变量中没有缓存
		if _, ok := actions[base]; !ok {
			actionKeys = append(actionKeys, path.Base(actionFile))
			actions[base] = actionFile
		}
	}
	//actionKeys保存了带数字的名件名，按数字名进行排序
	sort.Strings(actionKeys)

	var result = bytes.NewBufferString("")
	for _, key := range actionKeys {
		defer func() {
			if r := recover(); r != nil {
				err := fmt.Sprintf("Get error when render action %v: %v", key, r)
				utils.ExitWithMsg(err)
				result.WriteString(err)
			}
		}()

		base := regexp.MustCompile("^\\d+\\.").ReplaceAllString(key, "")
		if content, err := context.Asset(actions[base]); err == nil {
			if tmpl, err := template.New(filepath.Base(actions[base])).Funcs(context.FuncMap()).Parse(string(content)); err == nil {
				if err := tmpl.Execute(result, context); err != nil {
					result.WriteString(err.Error())
					utils.ExitWithMsg(err)
				}
			} else {
				result.WriteString(err.Error())
				utils.ExitWithMsg(err)
			}
		}
	}

	return template.HTML(strings.TrimSpace(result.String()))
}
func (context *Context) themesClass() (result string) {
	var results = map[string]bool{}
	if context.Resource != nil {
		for _, theme := range context.Resource.Config.Themes {
			if strings.HasPrefix(theme.GetName(), "-") {
				results[strings.TrimPrefix(theme.GetName(), "-")] = false
			} else if _, ok := results[theme.GetName()]; !ok {
				results[theme.GetName()] = true
			}
		}
	}

	var names []string
	for name, enabled := range results {
		if enabled {
			names = append(names, "qor-theme-"+name)
		}
	}
	return strings.Join(names, " ")
}

func (context *Context) logoutURL() string {
	if context.Admin.Auth != nil {
		return context.Admin.Auth.LogoutURL(context)
	}
	return ""
}

type menu struct {
	*Menu
	SubMenus []*menu
	Active   bool //当前request的url是否与此匹配
}

func (context *Context) getMenus() (menus []*menu) {
	var (
		globalMenu        = &menu{}
		mostMatchedMenu   *menu
		mostMatchedLength int
		addMenu           func(*menu, []*Menu)
	)

	addMenu = func(parent *menu, menus []*Menu) {
		for _, m := range menus {
			url := m.URL()

			if m.HasPermission(roles.Read, context.Context) {
				var menu = &menu{Menu: m}
				if strings.HasPrefix(context.Request.URL.Path, url) && len(url) > mostMatchedLength {
					mostMatchedMenu = menu
					mostMatchedLength = len(url)
				}

				addMenu(menu, menu.GetSubMenus())
				//menu必须有UR或者有子menu, 否则不显示
				if len(menu.SubMenus) > 0 || menu.URL() != "" {
					parent.SubMenus = append(parent.SubMenus, menu)
				}
			}
		}
	}
	//调用addMenu, 根据每一个menu的权限进行筛选
	addMenu(globalMenu, context.Admin.GetMenus())
	if context.Action != "search_center" && mostMatchedMenu != nil {
		mostMatchedMenu.Active = true
	}

	return globalMenu.SubMenus
}

func (context *Context) linkTo(text interface{}, link interface{}) template.HTML {
	text = reflect.Indirect(reflect.ValueOf(text)).Interface()
	if linkStr, ok := link.(string); ok {
		return template.HTML(fmt.Sprintf(`<a href="%v">%v</a>`, linkStr, text))
	}
	return template.HTML(fmt.Sprintf(`<a href="%v">%v</a>`, context.URLFor(link), text))
}

func (context *Context) t(values ...interface{}) template.HTML {
	switch len(values) {
	case 1:
		return context.Admin.T(context.Context, fmt.Sprint(values[0]), fmt.Sprint(values[0]))
	case 2:
		return context.Admin.T(context.Context, fmt.Sprint(values[0]), fmt.Sprint(values[1]))
	case 3:
		return context.Admin.T(context.Context, fmt.Sprint(values[0]), fmt.Sprint(values[1]), values[2:]...) //shared/sidebar.tmpl
	default:
		utils.ExitWithMsg("passed wrong params for T")
	}
	return ""
}
