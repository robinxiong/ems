package admin

import (
	"ems/core"
	"ems/core/resource"
	"ems/core/utils"
	"ems/roles"

	"github.com/jinzhu/gorm"
	"github.com/jinzhu/inflection"
	"fmt"
)

// Config resource config struct
type Config struct {
	Name       string   //指定Resource的名字字，默认为model struct的名字OrderItem，则为Order Item
	Menu       []string //父菜单的路径 "中国" "广东"
	Permission *roles.Permission
	Themes     []ThemeInterface
	Priority   int
	Singleton  bool //是否是当个Resource, 如果不是，则当前资源显示的菜单名称将为复数，如Colors, 同时在route.go中，注册resource路由时有用
	Invisible  bool //是否显示资源,默认值为false, 即显示资源, 则会将当前Resource添加到菜单
	PageCount  int
}

// Resource 在admin中是一个非常重要的概念，所以的数据库表model都抽像为一个资源, admin生成的管理接口，都是基于它的
type Resource struct {
	*resource.Resource //core/resource 用于将form数据写入到数据库
	Config             *Config
	ParentResource     *Resource  //可以参考route, RegisterRoute中，将使用到ParentResource
	SearchHandler      func(keyword string, context *core.Context) *gorm.DB
	params             string
	admin              *Admin //每一个Resource保留一份admin对像
	mounted            bool
}


func (res *Resource) GetAdmin() *Admin{
	return res.admin
}

// ToParam 将一个resource的名称,如Order Item转换为url param形式order_item，即替换空格为下划线
// 首先确定res.params是不为空，如果有res.params则直接返回， 否则调用model struct的ToParam(model struct实现了ToParam方法)
// 如果没有实现ToParam方法，则将res.Name转换成url允许的形式
// 用法可以参考AddResource中添加AddMenu (即生成菜单的RelativePath
// 第二个用法是route中RegisterRoute调用, 即向router中注册路径
func (res *Resource) ToParam() string {
	if res.params == "" {
		if value, ok := res.Value.(interface {
			ToParam() string
		}); ok {
			res.params = value.ToParam()
		} else {
			if res.Config.Singleton == true {
				res.params = utils.ToParamString(res.Name)
			} else {
				res.params = utils.ToParamString(inflection.Plural(res.Name))
			}
		}
	}
	return res.params
}
// ParamIDName ParamIDName return param name for primary key like :product_id
// 用法RegisterRoute
func (res Resource) ParamIDName() string {
	return fmt.Sprintf(":%v_id", inflection.Singular(utils.ToParamString(res.Name)))
}

//查看res model struct是否实现了resource.ConfigureResourceInterface接口，如果是则对它进行配置
func (res *Resource) configure() {
	//todo: 实现ConfigureResourceInterface
}
