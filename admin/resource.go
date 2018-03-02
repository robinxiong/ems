package admin

import (
	"ems/core/resource"
	"ems/roles"
	"ems/core"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/inflection"
	"ems/core/utils"
)

// Config resource config struct
type Config struct {
	Name       string
	Menu       []string
	Permission *roles.Permission
	Themes []ThemeInterface
	Priority   int
	Singleton  bool
	Invisible  bool
	PageCount  int
}

// Resource 在admin中是一个非常重要的概念，所以的数据库表model都抽像为一个资源, admin生成的管理接口，都是基于它的
type Resource struct {
	*resource.Resource
	Config *Config
	ParentResource *Resource
	SearchHandler func(keyword string, context *core.Context) *gorm.DB
	params   string
	admin *Admin
	mounted  bool
}



// ToParam used as urls to register routes for resource
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