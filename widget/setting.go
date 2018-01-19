package widget

import (
	"ems/serializable_meta"
	"time"
)

type QorWidgetSettingInterface interface {
	GetWidgetName() string
	SetWidgetName(string)
	GetGroupName() string
	SetGroupName(string)
	GetScope() string
	SetScope(string)
	GetTemplate() string
	SetTemplate(string)
	GetSourceType() string
	SetSourceType(string)
	GetSourceID() string
	SetSourceID(string)
	GetShared() bool
	SetShared(bool)
	serializable_meta.SerializableMetaInterface
}


// QorWidgetSetting default qor widget setting struct
type QorWidgetSetting struct {
	Name        string `gorm:"primary_key"`
	Scope       string `gorm:"primary_key;size:128;default:'default'"`
	SourceType  string `gorm:"primary_key;default:''"`
	SourceID    string `gorm:"primary_key;default:''"`
	Description string
	Shared      bool
	WidgetType  string
	GroupName   string
	Template    string
	serializable_meta.SerializableMeta
	CreatedAt time.Time
	UpdatedAt time.Time
}