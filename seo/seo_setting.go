package seo

import (
	"time"
	"ems/admin"
)

type SystemSEOSetting struct {
	Name        string `gorm:"primary_key"`
	Setting     Setting
	IsGlobalSEO bool

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`

	collection *Collection
}


type Setting struct {
	Title            string `gorm:"size:4294967295"`
	Description      string
	Keywords         string
	Type             string
	EnabledCustomize bool
	GlobalSetting    map[string]string
}



// Collection will hold registered seo configures and global setting definition and other configures
type Collection struct {
	Name            string
	SettingResource *admin.Resource

	registeredSEO  []*SEO
	resource       *admin.Resource
	globalResource *admin.Resource
	globalSetting  interface{}
}