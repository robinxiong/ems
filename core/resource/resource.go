package resource

import (
	"github.com/jinzhu/gorm"
	"ems/core"

	"ems/roles"
	"github.com/qor/qor"
)

type Resourcer interface{
	GetResource() *Resource
	GetMetas([]string) []Metaor
	CallFindMany(i interface{}, context *core.Context) error
	CallFindOne(i interface{}, value *MetaValue, context *core.Context) error
	CallSave(i interface{}, context *core.Context) error
	CallDelete(i interface{}, context *core.Context) error
	NewSlice() interface{}
	NewStruct() interface{}
}

// Resource用于定义资源的基础信息
type Resource struct {
	Name            string
	Value           interface{}
	PrimaryFields   []*gorm.StructField
	FindManyHandler func(interface{}, *core.Context) error
	FindOneHandler  func(interface{}, *MetaValues, *core.Context) error
	SaveHandler     func(interface{}, *core.Context) error
	DeleteHandler   func(interface{}, *core.Context) error
	Permission      *roles.Permission
	Validators      []*Validator
	Processors      []*Processor
	primaryField    *gorm.Field
}


type Validator struct {
	Name    string
	Handler func(interface{}, *MetaValues, *core.Context) error
}

// Processor processor struct
type Processor struct {
	Name    string
	Handler func(interface{}, *MetaValues, *qor.Context) error
}
