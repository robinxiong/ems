package admin

import (
	"ems/core/assetfs"
	"html/template"
)

var (
	globalViewPaths []string
	globalAssetFSes []assetfs.Interface
)



type I18n interface {
	Scope(scope string) I18n
	Default(value string) I18n
	T(locale string, key string, args ...interface{}) template.HTML
}