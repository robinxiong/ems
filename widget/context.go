package widget

import (
	"github.com/jinzhu/gorm"
	"html/template"
)

// Context widget context
type Context struct {
	Widgets          *Widgets
	DB               *gorm.DB
	AvailableWidgets []string
	Options          map[string]interface{}
	InlineEdit       bool
	SourceType       string
	SourceID         string
	FuncMaps         template.FuncMap
	WidgetSetting    QorWidgetSettingInterface
}