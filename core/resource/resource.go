package resource

import (
	"ems/core"

	"github.com/jinzhu/gorm"

	"ems/roles"

	"ems/core/utils"
	"fmt"
	"reflect"
)

type Resourcer interface {
	GetResource() *Resource
	GetMetas([]string) []Metaor
	CallFindMany(i interface{}, context *core.Context) error
	CallFindOne(i interface{}, value *MetaValues, context *core.Context) error
	CallSave(i interface{}, context *core.Context) error
	CallDelete(i interface{}, context *core.Context) error
	NewSlice() interface{}
	NewStruct() interface{}  //meta.Initialize()会调用它来创建一个新resource.Value的值
}

type ConfigureResourceBeforeInitializeInterface interface {
	ConfigureResourceBeforeInitialize(Resourcer)
}

//少部分model strcut实现了，e.g i18n, publish2, schedule, version
type ConfigureResourceInterface interface {
	ConfigureResource(Resourcer)
}

// Resource用于定义资源的基础信息, 它实际上一个数据库model的抽像
type Resource struct {
	Name            string
	Value           interface{} //所对应的model struct e.g&Resource{Value: &Product{}}
	PrimaryFields   []*gorm.StructField
	FindManyHandler func(interface{}, *core.Context) error //在crud.go中的CallFindMany中调用, 通常在resource.New()的时候使用默认的函数
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

	res.FindManyHandler = res.findManyHandler

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

func (res *Resource) GetMetas([]string) []Metaor {
	panic("implement me")
}

//初始化resource底层的model struct 数组（空数组) searcher.FindMany
func (res *Resource) NewSlice() interface{} {
	//如果没有绑定底层model struct
	if res.Value == nil {
		return nil
	}
	//res.Value类型为int, 则返回[]int
	sliceType := reflect.SliceOf(reflect.TypeOf(res.Value))
	//创建一个slice
	slice := reflect.MakeSlice(sliceType, 0, 0)

	slicePtr := reflect.New(sliceType) //指向到一个数组的指针
	slicePtr.Elem().Set(slice)         //设置一个

	return slicePtr.Interface() //将数组指针以interface{}返回

}

func (res *Resource) NewStruct() interface{} {
	if res.Value == nil {
		return nil
	}
	return reflect.New(utils.Indirect(reflect.ValueOf(res.Value)).Type()).Interface()
}

func (res *Resource) HasPermission(mode roles.PermissionMode, context *core.Context) bool {
	if res == nil || res.Permission == nil {
		return true
	}

	var roles = []interface{}{}
	//context.Roles是在admin/route.go ServeHTTP匹配到当前url，哪些角色
	for _, role := range context.Roles {
		roles = append(roles, role)
	}
	return res.Permission.HasPermission(mode, roles...)
}
