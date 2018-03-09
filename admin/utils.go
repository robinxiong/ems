package admin

import (
	"ems/core"
	"ems/core/assetfs"
	"html/template"
	"ems/roles"
)

type HasPermissioner interface {
	HasPermission(roles.PermissionMode, *core.Context) bool
}

//创建资源时model struct是否实现了这个接口，可以参考admin.newResource方法中的使用
type ResourceNamer interface {
	ResourceName() string
}

var (
	globalViewPaths []string
	globalAssetFSes []assetfs.Interface
)

type I18n interface {
	Scope(scope string) I18n
	Default(value string) I18n
	T(locale string, key string, args ...interface{}) template.HTML
}
