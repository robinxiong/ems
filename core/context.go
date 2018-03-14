package core

import (
	"net/http"
	"github.com/jinzhu/gorm"
)

type CurrentUser interface {
	DisplayName() string
}

type Context struct {
	Request *http.Request
	Writer http.ResponseWriter
	CurrentUser CurrentUser  //保存当前用户信息
	Roles []string
	DB *gorm.DB
	ResourceID string //admin/context.setResource
	Config *Config
	Errors
}


// Clone clone current context
func (context *Context) Clone() *Context {
	var clone = *context
	return &clone
}

// GetDB get db from current context
func (context *Context) GetDB() *gorm.DB {
	if context.DB != nil {
		return context.DB
	}
	return context.Config.DB
}

// SetDB set db into current context
func (context *Context) SetDB(db *gorm.DB) {
	context.DB = db
}