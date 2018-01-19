package resource

import (
	"ems/core"
	"ems/roles"
	"github.com/jinzhu/gorm"
)

//Metaor interface
type Metaor interface {
	GetName() string
	GetFieldName() string
	GetSetter() func(resource interface{}, metaValue *MetaValue, context *core.Context)
	GetFormattedValuer() func(resource interface{}, context core.Context) interface{}
	GetValuer() func(interface{}, *core.Context) interface{}
	GetResource() Resourcer
	GetMetas() []Metaor
	SetPermission(mode *roles.Permission)
	HasPermission(mode roles.PermissionMode, context *core.Context) bool
}


// ConfigureMetaBeforeInitializeInterface if a struct's field's type implemented this interface, it will be called when initializing a meta
type ConfigureMetaBeforeInitializeInterface interface {
	ConfigureMetaBeforeInitialize(Metaor)
}

// ConfigureMetaInterface if a struct's field's type implemented this interface, it will be called after configed
type ConfigureMetaInterface interface {
	ConfigureMeta(Metaor)
}

// MetaConfigInterface meta configuration interface
type MetaConfigInterface interface {
	ConfigureMetaInterface
}

// MetaConfig base meta config struct
type MetaConfig struct {
}

// ConfigureQorMeta implement the MetaConfigInterface
func (MetaConfig) ConfigureMeta(Metaor) {
}

type Meta struct {
	Name string
	FieldName string
	FieldStruct *gorm.StructField
	Setter func(resource interface{}, metaValue *MetaValue, context *core.Context)
	Valuer func(interface{}, *core.Context) interface{}
	FormatterValuer func(interface{}, *core.Context) interface{}
	Config MetaConfigInterface
	BaseResource Resourcer
	Permission *roles.Permission
}