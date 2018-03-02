package admin

import (
	"ems/core/assetfs"
	"ems/core/utils"
	"path/filepath"

	"ems/core"

	"html/template"

	"github.com/jinzhu/gorm"
	"github.com/theplant/cldr"
)

type AdminConfig struct {
	SiteName string //用于layout.tmpl的title, 而且javascript, 以及css都依赖这个值
	DB       *gorm.DB
	AssetFS  assetfs.Interface //core.assetfs 用来获取template, css, js等文件，assetFS中会绑定这些文件所在的目录
	I18n     I18n              //context中，用来翻译的，比如pageTitle func map
	Auth     Auth    //用于认证的接口，包含getCurrentUser, LoginURL, logoutURL等方法
}

type Admin struct {
	*AdminConfig
	router *Router
}

func New(config interface{}) *Admin {
	admin := Admin{
		router: newRouter(),
	}

	if c, ok := config.(*core.Config); ok {
		admin.AdminConfig = &AdminConfig{DB: c.DB}
	} else if c, ok := config.(*AdminConfig); ok {
		admin.AdminConfig = c
	} else {
		admin.AdminConfig = &AdminConfig{}
	}

	if admin.AssetFS == nil {
		admin.AssetFS = assetfs.AssetFS().NameSpace("admin") //global assetfs,创建一个admin的命名
	}
	//设置路径，这样在controller中查找layout.tmpl, 图片都资源，依赖于assetFS中的path
	admin.SetAssetFS(admin.AssetFS)

	return &admin
}

func (admin *Admin) GetRouter() *Router {
	return admin.router
}

//设置admin的asset查找路径
func (admin *Admin) SetAssetFS(assetFS assetfs.Interface) {
	admin.AssetFS = assetFS
	globalAssetFSes = append(globalAssetFSes, assetFS)
	//注册当前网站的asset目录 /site/app/views/, 它可以替换admin/views/中的模板
	admin.AssetFS.RegisterPath(filepath.Join(utils.AppRoot, "app/views/system"))
	admin.RegisterViewPath("ems/admin/views")
	//在vendor和gopath中搜索viewPath
	for _, viewPath := range globalViewPaths {
		admin.RegisterViewPath(viewPath)
	}
}

func (admin *Admin) RegisterViewPath(pth string) {
	if admin.AssetFS.RegisterPath(filepath.Join(utils.AppRoot, "vendor", pth)) != nil {
		for _, gopath := range utils.GOPATH() {
			if admin.AssetFS.RegisterPath(filepath.Join(gopath, "src", pth)) == nil {
				break
			}
		}
	}
}

func (admin *Admin) T(context *core.Context, key string, value string, values ...interface{}) template.HTML {
	//如果admin中没有配置I18n(admin/utils.go)
	locale := utils.GetLocale(context)

	if admin.I18n == nil {
		if result, err := cldr.Parse(locale, value, values...); err == nil {
			return template.HTML(result)
		}
		return template.HTML(key)
	}
	//通常调用，除非实现了admin.I18n接口
	return admin.I18n.Default(value).T(locale, key, values...)
}


