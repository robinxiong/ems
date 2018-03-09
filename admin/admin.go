package admin

import (
	"ems/core/assetfs"
	"ems/core/utils"
	"path/filepath"

	"ems/core"

	"html/template"

	"ems/core/resource"
	"reflect"

	"github.com/jinzhu/gorm"
	"github.com/theplant/cldr"
	"github.com/jinzhu/inflection"
)

type AdminConfig struct {
	SiteName string //用于layout.tmpl的title, 而且javascript, 以及css都依赖这个值
	DB       *gorm.DB
	AssetFS  assetfs.Interface //core.assetfs 用来获取template, css, js等文件，assetFS中会绑定这些文件所在的目录
	I18n     I18n              //context中，用来翻译的，比如pageTitle func map
	Auth     Auth              //用于认证的接口，包含getCurrentUser, LoginURL, logoutURL等方法
}

type Admin struct {
	*AdminConfig
	router *Router
	menus  []*Menu //保存后台使用到的菜单，它通过main.go中设置, 比如site/app/admin/dashboard.go或者/products.go
	resources        []*Resource //保存所有创建的资源, AddResource时添加资源
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

//Resource

//AddResource的第一个参数是model的值, config 为resource.go/Config 配置名称和menu，以及priority
func (admin *Admin) AddResource(value interface{}, config ...*Config) *Resource {
	res := admin.newResource(value, config...)
	admin.resources = append(admin.resources, res)
	//如果resource model struct实现了resource.ConfigureResourceInterface，则对model进行配置
	res.configure()

	if !res.Config.Invisible {
		menuName := res.Name
		if !res.Config.Singleton {
			menuName = inflection.Plural(res.Name)
		}

		admin.AddMenu(&Menu{Name:menuName, Priority:res.Config.Priority, Ancestors: res.Config.Menu, RelativePath:res.ToParam()})

		//添加完菜单后，注册路由
		admin.RegisterResourceRouters(res, "create", "update", "read", "delete")
	}


	return res
}

func (admin *Admin) newResource(value interface{}, config ...*Config) *Resource {
	var configuration *Config
	if len(config) > 0 {
		configuration = config[0]
	}

	if configuration == nil {
		configuration = &Config{}
	}
	res := &Resource{
		Resource: resource.New(value), //core.Resource的New方法， resource的Name会调用core/utils.HumanizeString "OrderItem" -> "Order Item"
		Config:   configuration,
		admin:    admin,
	}

	//将config中的permission复制到core.Resource
	res.Permission = configuration.Permission

	//资源的名称, 如果定义了configuration.Name为空，或者model设置了ResourceName方法则调用, 默认为model struct的名字
	if configuration.Name != "" {
		res.Name = configuration.Name
	} else if namer, ok := value.(ResourceNamer); ok {
		res.Name = namer.ResourceName()
	}

	//model struct是否嵌入其它struct, 比如publish2/publish2 Publish struct
	modelType := utils.ModelType(res.Value)
	for i := 0; i < modelType.NumField(); i++ {
		//字段是否实现了resource.ConfigureResourceBeforeInitializeInterface
		if fieldStruct := modelType.Field(i); fieldStruct.Anonymous {
			if injector, ok := reflect.New(fieldStruct.Type).Interface().(resource.ConfigureResourceBeforeInitializeInterface); ok {
				injector.ConfigureResourceBeforeInitialize(res)
			}
		}
	}
	if injector, ok := res.Value.(resource.ConfigureResourceBeforeInitializeInterface); ok {
		injector.ConfigureResourceBeforeInitialize(res)
	}


	return res
}
