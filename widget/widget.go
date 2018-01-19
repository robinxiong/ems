package widget

import (
	"html/template"
	"ems/admin"
	"ems/roles"
	"github.com/jinzhu/gorm"
	"ems/assetfs"
)

// Config widget config
type Config struct {
	DB            *gorm.DB
	Admin         *admin.Admin
	PreviewAssets []string
}



type Widgets struct {
	funcMaps              template.FuncMap
	Config                *Config
	Resource              *admin.Resource
	AssetFS               assetfs.Interface
	WidgetSettingResource *admin.Resource
}

type Widget struct {
	Name          string
	PreviewIcon   string
	Group         string
	Templates     []string
	Setting       *admin.Resource
	Permission    *roles.Permission
	InlineEditURL func(*Context) string
	Context       func(context *Context, setting interface{}) *Context
}