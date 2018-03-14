package admin

import (
	"bytes"
	"ems/core/utils"
	"ems/roles"
	"fmt"
	"html/template"
	"math/rand"
	"net/url"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"database/sql/driver"
	"ems/core"
	"runtime/debug"

	"github.com/jinzhu/gorm"
)

//for context

func (context *Context) FuncMap() template.FuncMap {
	funcMap := template.FuncMap{
		"page_title":             context.pageTitle,
		"t":                      context.t,                    //翻译成本地语言，它调用admin.T
		"stylesheet_tag":         context.styleSheetTag,        //根据指定的字符串，返回template.HTML, 比如admin_default, 它返回/admin/assets/stylesheets/admin_default.css, 对于assets的请求，可以查看admin.Route中的serveHTTP
		"load_admin_stylesheets": context.loadAdminStyleSheets, //返回指定siteName的css, 如果没有指定siteName, 则返回application_css
		"load_theme_stylesheets": context.loadThemeStyleSheets,
		"javascript_tag":         context.javaScriptTag, //从/admin/views/assets/javascripts中加载指定名字的js文件
		"load_actions":           context.loadActions,   //加载views/action下的子模板, action会按照文件前面的数字进行排名
		"qor_theme_class":        context.themesClass,   //返回theme style的类名称，用于body标签
		"render":                 context.Render,        //读取指定的模板
		"logout_url":             context.logoutURL,     //sidebar.tmpl获取logout url
		"get_menus":              context.getMenus,      //获取系统菜单，并且传递给sidebar.tmpl
		"url_for":                context.URLFor,
		"link_to":                context.linkTo,               //翻译link_to的名称
		"load_admin_javascripts": context.loadAdminJavaScripts, //根据site来加载js文件

		"primary_key_of": context.primaryKeyOf,
		"unique_key_of":  context.uniqueKeyOf,
		//view/index/table.tmpl
		"convert_sections_to_metas": context.convertSectionToMetas,
		"index_sections":            context.indexSections,
		"is_sortable_meta":          context.isSortableMeta,

		//meta
		"meta_label": func(meta *Meta) template.HTML {
			key := fmt.Sprintf("%v.attributes.%v", meta.baseResource.ToParam(), meta.Label)
			return context.Admin.T(context.Context, key, meta.Label)
		},
		"render_meta": func(value interface{}, meta *Meta, types ...string) template.HTML {
			var (
				result = bytes.NewBufferString("")
				typ    = "index"
			)

			for _, t := range types {
				typ = t
			}

			context.renderMeta(meta, value, []string{}, typ, result)
			return template.HTML(result.String())
		},
	}
	return funcMap
}
func (context *Context) URLFor(value interface{}, resources ...*Resource) string {
	getPrefix := func(res *Resource) string {
		var params string
		for res.ParentResource != nil {
			params = path.Join(res.ParentResource.ToParam(), res.ParentResource.GetPrimaryValue(context.Request), params)
			res = res.ParentResource
		}
		return path.Join(res.GetAdmin().router.Prefix, params)
	}

	if admin, ok := value.(*Admin); ok {
		return admin.router.Prefix
	} else if res, ok := value.(*Resource); ok {
		return path.Join(getPrefix(res), res.ToParam())
	} else {
		var res *Resource

		if len(resources) > 0 {
			res = resources[0]
		}

		if res == nil {
			res = context.Admin.GetResource(reflect.Indirect(reflect.ValueOf(value)).Type().String())
		}

		if res != nil {
			if res.Config.Singleton {
				return path.Join(getPrefix(res), res.ToParam())
			}

			var (
				scope         = context.GetDB().NewScope(value)
				primaryFields []string
				primaryValues = map[string]string{}
			)

			for _, primaryField := range res.PrimaryFields {
				if field, ok := scope.FieldByName(primaryField.Name); ok {
					primaryFields = append(primaryFields, fmt.Sprint(field.Field.Interface())) // TODO improve me
				}
			}

			for _, field := range scope.PrimaryFields() {
				useAsPrimaryField := false
				for _, primaryField := range res.PrimaryFields {
					if field.DBName == primaryField.DBName {
						useAsPrimaryField = true
						break
					}
				}

				if !useAsPrimaryField {
					primaryValues[fmt.Sprintf("primary_key[%v_%v]", scope.TableName(), field.DBName)] = fmt.Sprint(reflect.Indirect(field.Field).Interface())
				}
			}

			urlPath := path.Join(getPrefix(res), res.ToParam(), strings.Join(primaryFields, ","))

			if len(primaryValues) > 0 {
				var primaryValueParams []string
				for key, value := range primaryValues {
					primaryValueParams = append(primaryValueParams, fmt.Sprintf("%v=%v", key, url.QueryEscape(value)))
				}
				urlPath = urlPath + "?" + strings.Join(primaryValueParams, "&")
			}
			return urlPath
		}
	}
	return ""
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
func (context *Context) getResource(resources ...*Resource) *Resource {
	for _, res := range resources {
		return res
	}
	return context.Resource
}

func (context *Context) primaryKeyOf(value interface{}) interface{} {
	if reflect.Indirect(reflect.ValueOf(value)).Kind() == reflect.Struct {
		scope := &gorm.Scope{Value: value}
		return fmt.Sprint(scope.PrimaryKeyValue())
	}
	return fmt.Sprint(value)
}

func (context *Context) uniqueKeyOf(value interface{}) interface{} {
	if reflect.Indirect(reflect.ValueOf(value)).Kind() == reflect.Struct {
		scope := &gorm.Scope{Value: value}
		var primaryValues []string
		for _, primaryField := range scope.PrimaryFields() {
			primaryValues = append(primaryValues, fmt.Sprint(primaryField.Field.Interface()))
		}
		primaryValues = append(primaryValues, fmt.Sprint(rand.Intn(1000)))
		return utils.ToParamString(url.QueryEscape(strings.Join(primaryValues, "_")))
	}
	return fmt.Sprint(value)
}

//meta
func (context *Context) hasCreatePermission(permissioner HasPermissioner) bool {
	return permissioner.HasPermission(roles.Create, context.Context)
}

func (context *Context) hasReadPermission(permissioner HasPermissioner) bool {
	return permissioner.HasPermission(roles.Read, context.Context)
}

func (context *Context) hasUpdatePermission(permissioner HasPermissioner) bool {
	return permissioner.HasPermission(roles.Update, context.Context)
}

func (context *Context) hasDeletePermission(permissioner HasPermissioner) bool {
	return permissioner.HasPermission(roles.Delete, context.Context)
}

func (context *Context) hasChangePermission(permissioner HasPermissioner) bool {
	if context.Action == "new" {
		return context.hasCreatePermission(permissioner)
	}
	return context.hasUpdatePermission(permissioner)
}

func (context *Context) isNewRecord(value interface{}) bool {
	if value == nil {
		return true
	}
	return context.GetDB().NewRecord(value)
}
func (context *Context) valueOf(valuer func(interface{}, *core.Context) interface{}, value interface{}, meta *Meta) interface{} {
	if valuer != nil {
		reflectValue := reflect.ValueOf(value)
		if reflectValue.Kind() != reflect.Ptr {
			reflectPtr := reflect.New(reflectValue.Type())
			reflectPtr.Elem().Set(reflectValue)
			value = reflectPtr.Interface()
		}

		result := valuer(value, context.Context)

		if reflectValue := reflect.ValueOf(result); reflectValue.IsValid() {
			if reflectValue.Kind() == reflect.Ptr {
				if reflectValue.IsNil() || !reflectValue.Elem().IsValid() {
					return nil
				}

				result = reflectValue.Elem().Interface()
			}

			if meta.Type == "number" || meta.Type == "float" {
				if context.isNewRecord(value) && equal(reflect.Zero(reflect.TypeOf(result)).Interface(), result) {
					return nil
				}
			}
			return result
		}
		return nil
	}

	utils.ExitWithMsg(fmt.Sprintf("No valuer found for meta %v of resource %v", meta.Name, meta.baseResource.Name))
	return nil
}

// FormattedValueOf return formatted value of a meta for current resource
func (context *Context) FormattedValueOf(value interface{}, meta *Meta) interface{} {
	result := context.valueOf(meta.GetFormattedValuer(), value, meta)
	if resultValuer, ok := result.(driver.Valuer); ok {
		if result, err := resultValuer.Value(); err == nil {
			return result
		}
	}

	return result
}
func (context *Context) renderSections(value interface{}, sections []*Section, prefix []string, writer *bytes.Buffer, kind string) {
	for _, section := range sections {
		var rows []struct {
			Length      int
			ColumnsHTML template.HTML
		}

		for _, column := range section.Rows {
			columnsHTML := bytes.NewBufferString("")
			for _, col := range column {
				meta := section.Resource.GetMeta(col)
				if meta != nil {
					context.renderMeta(meta, value, prefix, kind, columnsHTML)
				}
			}

			rows = append(rows, struct {
				Length      int
				ColumnsHTML template.HTML
			}{
				Length:      len(column),
				ColumnsHTML: template.HTML(string(columnsHTML.Bytes())),
			})
		}

		if len(rows) > 0 {
			var data = map[string]interface{}{
				"Section": section,
				"Title":   template.HTML(section.Title),
				"Rows":    rows,
			}
			if content, err := context.Asset("metas/section.tmpl"); err == nil {
				if tmpl, err := template.New("section").Funcs(context.FuncMap()).Parse(string(content)); err == nil {
					tmpl.Execute(writer, data)
				}
			}
		}
	}
}
func (context *Context) renderMeta(meta *Meta, value interface{}, prefix []string, metaType string, writer *bytes.Buffer) {
	var (
		err      error
		funcsMap = context.FuncMap()
	)
	prefix = append(prefix, meta.Name)

	var generateNestedRenderSections = func(kind string) func(interface{}, []*Section, int) template.HTML {
		return func(value interface{}, sections []*Section, index int) template.HTML {
			var result = bytes.NewBufferString("")
			var newPrefix = append([]string{}, prefix...)

			if index >= 0 {
				last := newPrefix[len(newPrefix)-1]
				newPrefix = append(newPrefix[:len(newPrefix)-1], fmt.Sprintf("%v[%v]", last, index))
			}

			if len(sections) > 0 {
				for _, field := range context.GetDB().NewScope(value).PrimaryFields() {
					if meta := sections[0].Resource.GetMeta(field.Name); meta != nil {
						context.renderMeta(meta, value, newPrefix, kind, result)
					}
				}

				context.renderSections(value, sections, newPrefix, result, kind)
			}

			return template.HTML(result.String())
		}
	}

	funcsMap["has_change_permission"] = func(permissioner HasPermissioner) bool {
		if context.GetDB().NewScope(value).PrimaryKeyZero() {
			return context.hasCreatePermission(permissioner)
		}
		return context.hasUpdatePermission(permissioner)
	}
	funcsMap["render_nested_form"] = generateNestedRenderSections("form")

	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			writer.WriteString(fmt.Sprintf("Get error when render template for meta %v (%v): %v", meta.Name, meta.Type, r))
		}
	}()

	var (
		tmpl    = template.New(meta.Type + ".tmpl").Funcs(funcsMap)
		content []byte
	)

	switch {
	case meta.Config != nil:
		if templater, ok := meta.Config.(interface {
			GetTemplate(context *Context, metaType string) ([]byte, error)
		}); ok {
			if content, err = templater.GetTemplate(context, metaType); err == nil {
				tmpl, err = tmpl.Parse(string(content))
				break
			}
		}
		fallthrough
	default:
		if content, err = context.Asset(fmt.Sprintf("%v/metas/%v/%v.tmpl", meta.baseResource.ToParam(), metaType, meta.Name), fmt.Sprintf("metas/%v/%v.tmpl", metaType, meta.Type)); err == nil {
			tmpl, err = tmpl.Parse(string(content))
		} else if metaType == "index" {
			tmpl, err = tmpl.Parse("{{.Value}}")
		} else {
			err = fmt.Errorf("haven't found %v template for meta %v", metaType, meta.Name)
		}
	}

	if err == nil {
		var scope = context.GetDB().NewScope(value)
		var data = map[string]interface{}{
			"Context":       context,
			"BaseResource":  meta.baseResource,
			"Meta":          meta,
			"ResourceValue": value,
			"Value":         context.FormattedValueOf(value, meta),
			"Label":         meta.Label,
			"InputName":     strings.Join(prefix, "."),
		}

		if !scope.PrimaryKeyZero() {
			data["InputId"] = utils.ToParamString(fmt.Sprintf("%v_%v_%v", scope.GetModelStruct().ModelType.Name(), scope.PrimaryKeyValue(), meta.Name))
		}

		data["CollectionValue"] = func() [][]string {
			fmt.Printf("%v: Call .CollectionValue from views already Deprecated, get the value with `.Meta.Config.GetCollection .ResourceValue .Context`", meta.Name)
			return meta.Config.(interface {
				GetCollection(value interface{}, context *Context) [][]string
			}).GetCollection(value, context)
		}

		err = tmpl.Execute(writer, data)
	}

	if err != nil {
		msg := fmt.Sprintf("got error when render %v template for %v(%v): %v", metaType, meta.Name, meta.Type, err)
		fmt.Fprint(writer, msg)
		utils.ExitWithMsg(msg)
	}
}

//view/index/table.tmpl
//如果参数为空，则返回context.Resource
func (context *Context) indexSections(resources ...*Resource) []*Section {
	res := context.getResource(resources...)
	return res.allowedSections(res.IndexAttrs(), context, roles.Read)
}

func (context *Context) convertSectionToMetas(res *Resource, sections []*Section) []*Meta {

	return res.ConvertSectionToMetas(sections)
}

func (context *Context) isSortableMeta(meta *Meta) bool {
	for _, attr := range context.Resource.SortableAttrs() {
		if attr == meta.Name && meta.FieldStruct != nil && meta.FieldStruct.IsNormal && meta.FieldStruct.DBName != "" {
			return true
		}
	}
	return false
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
