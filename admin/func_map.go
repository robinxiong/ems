package admin

import (
	"ems/core/utils"
	"fmt"
	"html/template"
	"path"
	"strings"
	"path/filepath"
	"regexp"
	"bytes"
	"sort"
)

//for context

func (context *Context) FuncMap() template.FuncMap {
	funcMap := template.FuncMap{
		"page_title":     context.pageTitle,
		"t":              context.t,             //翻译成本地语言，它调用admin.T
		"stylesheet_tag": context.styleSheetTag, //根据指定的字符串，返回template.HTML, 比如admin_default, 它返回/admin/assets/stylesheets/admin_default.css, 对于assets的请求，可以查看admin.Route中的serveHTTP
		"load_admin_stylesheets": context.loadAdminStyleSheets, //返回指定siteName的css, 如果没有指定siteName, 则返回application_css
		"load_theme_stylesheets": context.loadThemeStyleSheets,
		"javascript_tag":         context.javaScriptTag,
		"load_actions":  context.loadActions,  //加载views/action下的子模板
	}
	return funcMap
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
func (context *Context) loadAdminStyleSheets() template.HTML{
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
func (context *Context) loadActions(action string) template.HTML{
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


func (context *Context) t(values ...interface{}) template.HTML {
	switch len(values) {
	case 1:
		return context.Admin.T(context.Context, fmt.Sprint(values[0]), fmt.Sprint(values[0]))
	case 2:
		return context.Admin.T(context.Context, fmt.Sprint(values[0]), fmt.Sprint(values[1]))
	case 3:
		return context.Admin.T(context.Context, fmt.Sprint(values[0]), fmt.Sprint(values[1]), values[2:]...)
	default:
		utils.ExitWithMsg("passed wrong params for T")
	}
	return ""
}
