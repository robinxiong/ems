package resource

import (
	"ems/core"

	"github.com/jinzhu/gorm"

	"ems/roles"

	"ems/core/utils"
	"fmt"
)

type Resourcer interface {
	GetResource() *Resource
	GetMetas([]string) []Metaor
	CallFindMany(i interface{}, context *core.Context) error
	CallFindOne(i interface{}, value *MetaValue, context *core.Context) error
	CallSave(i interface{}, context *core.Context) error
	CallDelete(i interface{}, context *core.Context) error
	NewSlice() interface{}
	NewStruct() interface{}
}

// Resource用于定义资源的基础信息, 它实际上一个数据库model的抽像
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

//初始化一个resource
func New(value interface{}) *Resource {
	var (
		//解析model value, 获取数据库表名
		name = utils.HumanizeString(utils.ModelType(value).Name())
		res  = &Resource{Value: value, Name: name} //基于model, 创建一个资源
	)

	return res
}

// GetResource return itself to match interface `Resourcer`
func (res *Resource) GetResource() *Resource {
	return res
}

// SetPrimaryFields set primary fields
func (res *Resource) SetPrimaryFields(fields ...string) error {
	scope := gorm.Scope{Value: res.Value}
	res.PrimaryFields = nil

	if len(fields) > 0 {
		for _, fieldName := range fields {
			if field, ok := scope.FieldByName(fieldName); ok {
				res.PrimaryFields = append(res.PrimaryFields, field.StructField)
			} else {
				return fmt.Errorf("%v is not a valid field for resource %v", fieldName, res.Name)
			}
		}
		return nil
	}

	if primaryField := scope.PrimaryField(); primaryField != nil {
		res.PrimaryFields = []*gorm.StructField{primaryField.StructField}
		return nil
	}

	return fmt.Errorf("no valid primary field for resource %v", res.Name)
}

type Validator struct {
	Name    string
	Handler func(interface{}, *MetaValues, *core.Context) error
}

// Processor processor struct
type Processor struct {
	Name    string
	Handler func(interface{}, *MetaValues, *core.Context) error
}
