package admin

import (
	"github.com/jinzhu/gorm"
	"ems/core"
)

type Scope struct {
	Name string
	Label string
	Group string
	Visible func(context *Context)bool
	Handler func(*gorm.DB,*core.Context) *gorm.DB
	Default bool
}

func (res *Resource) Scope(scope *Scope) {
	if scope.Label == "" {
		scope.Label = scope.Name
	}
}