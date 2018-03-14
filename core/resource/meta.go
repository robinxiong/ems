package resource

import (
	"database/sql"
	"ems/core"
	"ems/core/utils"
	"ems/roles"
	"ems/validations"
	"fmt"
	"reflect"
	"runtime/debug"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

//Metaor interface
type Metaor interface {
	GetName() string
	GetFieldName() string
	//返回Meta的Setter
	GetSetter() func(resource interface{}, metaValue *MetaValue, context *core.Context)
	GetFormattedValuer() func(resource interface{}, context *core.Context) interface{}
	//返回Meta对的Valuer，Valuer是在Initialize中设置的
	GetValuer() func(interface{}, *core.Context) interface{}
	GetResource() Resourcer
	GetMetas() []Metaor
	SetPermission(mode *roles.Permission)
	HasPermission(mode roles.PermissionMode, context *core.Context) bool
}

// ConfigureMetaBeforeInitializeInterface if a struct's field's type implemented this interface, it will be called when initializing a meta
// admin/meta_test.go 调用media/base.go
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

	FieldName   string            //所对应model struct的字段名称, 通常跟Name相同
	FieldStruct *gorm.StructField //meta所对应字段的gorm.StructField
	//设置meta.Setter为一个函数，第一个参娄秋数据库记录，第二个为MetaValue, 第二个参数为context
	//它首先根据meta.FieldName来查找record中的字段的值, 然后调用setter方法，将metaValue中的值，赋值给field
	Setter func(resource interface{}, metaValue *MetaValue, context *core.Context)
	//从model struct中，茆以meta所对应字段的值, 这个函数通常在Initialize中设置
	Valuer          func(interface{}, *core.Context) interface{} //meta.Valuer的第一个参数为model struct的记录值，第二个为context, 返回meta.FieldName中字段的值
	FormattedValuer func(interface{}, *core.Context) interface{}
	Config          MetaConfigInterface
	BaseResource    Resourcer   //所属的Resource, 当前meta的父Resource
	Permission      *roles.Permission
	Resource        Resourcer  //meta的字段如果为struct, 在baseResource的基础上在创建一个subResource
}

func (meta Meta) GetName() string {
	panic("implement me")
}

func (meta Meta) GetSetter() func(resource interface{}, metaValue *MetaValue, context *core.Context) {
	return meta.Setter
}

func (meta Meta) GetFormattedValuer() func(resource interface{}, context *core.Context) interface{} {
	if meta.FormattedValuer != nil {
		return meta.FormattedValuer
	}
	return meta.Valuer
}

func (meta Meta) GetValuer() func(interface{}, *core.Context) interface{} {
	return meta.Valuer
}

func (meta Meta) GetResource() Resourcer {
	panic("implement me")
}

func (meta Meta) GetMetas() []Metaor {
	panic("implement me")
}

func (meta Meta) SetPermission(mode *roles.Permission) {
	panic("implement me")
}

func (meta Meta) HasPermission(mode roles.PermissionMode, context *core.Context) bool {
	if meta.Permission == nil {
		return true
	}
	var roles = []interface{}{}
	for _, role := range context.Roles {
		roles = append(roles, role)
	}
	return meta.Permission.HasPermission(mode, roles...)
}

//返回meta所对应数据库字段的名称
func (meta Meta) GetFieldName() string {
	return meta.FieldName
}

// GetBaseResource get base resource from meta
func (meta Meta) GetBaseResource() Resourcer {
	return meta.BaseResource
}

// PreInitialize when will be run before initialize, used to fill some basic necessary information
// 主要是初始化meta.FieldName和meta.FieldStruct, 它是meta所对应的model struct字段的gorm.StructField
// FieldStruct将在setupValuer中使用
func (meta *Meta) PreInitialize() error {
	if meta.Name == "" {
		utils.ExitWithMsg("Meta should have name: %v", reflect.TypeOf(meta))
	} else if meta.FieldName == "" {
		meta.FieldName = meta.Name
	}

	// parseNestedField used to handle case like Profile.Name  (User下面内嵌了Profile, 而meta指定的名称为Profile.Name
	var parseNestedField = func(value reflect.Value, name string) (reflect.Value, string) {
		fields := strings.Split(name, ".")
		value = reflect.Indirect(value)
		//一级一级的获取到每一个字段的实际struct, Name之前的那个struct
		for _, field := range fields[:len(fields)-1] {
			value = value.FieldByName(field)
		}
		//返回这个struct和字段的名称
		return value, fields[len(fields)-1]
	}

	//第一个参数通常是能过scope.GetStructFields获取到一个model struct的所有字段信息，然后跟第二个参数匹配
	//返回匹配到的字段
	var getField = func(fields []*gorm.StructField, name string) *gorm.StructField {
		for _, field := range fields {
			if field.Name == name || field.DBName == name {
				return field
			}
		}
		return nil
	}
	// meta_test.go TestNestedField
	var nestedField = strings.Contains(meta.FieldName, ".")
	var scope = &gorm.Scope{Value: meta.BaseResource.GetResource().Value}
	if nestedField {
		subModel, name := parseNestedField(reflect.ValueOf(meta.BaseResource.GetResource().Value), meta.FieldName)
		meta.FieldStruct = getField(scope.New(subModel.Interface()).GetStructFields(), name)
	} else {
		meta.FieldStruct = getField(scope.GetStructFields(), meta.FieldName)
	}
	return nil
}

// Initialize initialize meta, will set valuer, setter if haven't configure it
// 它在meta.configure的时候调用
func (meta *Meta) Initialize() error {
	// Set Valuer for Meta， 调用默认的
	if meta.Valuer == nil {

		setupValuer(meta, meta.FieldName, meta.GetBaseResource().NewStruct())
	}

	if meta.Valuer == nil {
		utils.ExitWithMsg("Meta %v is not supported for resource %v, no `Valuer` configured for it", meta.FieldName, reflect.TypeOf(meta.BaseResource.GetResource().Value))
	}

	// Set Setter for Meta
	if meta.Setter == nil {
		setupSetter(meta, meta.FieldName, meta.GetBaseResource().NewStruct())
	}
	return nil
}

//设置meta.Valuer的值为一个函数,用于获取meta的值
//meta.Valuer的第一个参数为model struct的记录值，第二个为context, 返回meta.FieldName中字段的值
func setupValuer(meta *Meta, fieldName string, record interface{}) {
	nestedField := strings.Contains(fieldName, ".")

	// Setup nested fields
	if nestedField {
		fieldNames := strings.Split(fieldName, ".")

		//getNetstedModel的第一个参数是当前model的值
		//第二个参数为 x.y的形式User下面的Profile.Name, 则为Profile.Name, getNestedModel将返回第二个参数中的第一个即Profile类型
		// admin/meta_test.go TestNestedField
		setupValuer(meta, strings.Join(fieldNames[1:], "."), getNestedModel(record, strings.Join(fieldNames[0:2], "."), nil))

		oldValuer := meta.Valuer
		meta.Valuer = func(record interface{}, context *core.Context) interface{} {
			return oldValuer(getNestedModel(record, strings.Join(fieldNames[0:2], "."), context), context)
		}
		return
	}

	//PreInitialize中初始化了FieldStruct, 即它段的qorm.StructField

	if meta.FieldStruct != nil {
		meta.Valuer = func(value interface{}, context *core.Context) interface{} {
			scope := context.GetDB().NewScope(value)

			if f, ok := scope.FieldByName(fieldName); ok {
				if relationship := f.Relationship; relationship != nil && f.Field.CanAddr() && !scope.PrimaryKeyZero() {
					if (relationship.Kind == "has_many" || relationship.Kind == "many_to_many") && f.Field.Len() == 0 {
						// meta_test.go /TestGetSliceMetaValue
						context.GetDB().Model(value).Related(f.Field.Addr().Interface(), fieldName)

					} else if (relationship.Kind == "has_one" || relationship.Kind == "belongs_to") && context.GetDB().NewScope(f.Field.Interface()).PrimaryKeyZero() {
						// meta_test.go /TestGetStructMetaValue
						if f.Field.Kind() == reflect.Ptr && f.Field.IsNil() {
							f.Field.Set(reflect.New(f.Field.Type().Elem()))
						}

						context.GetDB().Model(value).Related(f.Field.Addr().Interface(), fieldName)
					}
				}

				return f.Field.Interface()
			}

			return ""
		}
	}
}

//设置meta.Setter为一个函数，第一个参娄秋数据库记录，第二个为MetaValue, 第二个参数为context
//它首先根据meta.FieldName来查找record中的字段的值, 然后调用setter方法，将metaValue中的值，赋值给field
//admin/meta_test.go TestStringMetaSetter
func setupSetter(meta *Meta, fieldName string, record interface{}) {
	nestedField := strings.Contains(fieldName, ".")

	//如果是嵌套的字段，找到子字段, User下内嵌一个Profile, meta的名称为Profile.Name
	if nestedField {
		fieldNames := strings.Split(fieldName, ".")

		setupSetter(meta, strings.Join(fieldNames[1:], "."), getNestedModel(record, strings.Join(fieldNames[0:2], "."), nil))

		oldSetter := meta.Setter
		meta.Setter = func(record interface{}, metaValue *MetaValue, context *core.Context) {
			oldSetter(getNestedModel(record, strings.Join(fieldNames[0:2], "."), context), metaValue, context)
		}
		return
	}

	//定义一个高阶函数, 它接收一个setter,
	commonSetter := func(setter func(field reflect.Value, metaValue *MetaValue, context *core.Context, record interface{})) func(record interface{}, metaValue *MetaValue, context *core.Context) {
		return func(record interface{}, metaValue *MetaValue, context *core.Context) {
			if metaValue == nil {
				return
			}

			defer func() {
				if r := recover(); r != nil {
					debug.PrintStack()
					context.AddError(validations.NewError(record, meta.Name, fmt.Sprintf("Failed to set Meta %v's value with %v, got %v", meta.Name, metaValue.Value, r)))
				}
			}()

			field := utils.Indirect(reflect.ValueOf(record)).FieldByName(fieldName)
			if field.Kind() == reflect.Ptr {
				if field.IsNil() && utils.ToString(metaValue.Value) != "" {
					field.Set(utils.NewValue(field.Type()).Elem())
				}

				if utils.ToString(metaValue.Value) == "" {
					field.Set(reflect.Zero(field.Type()))
					return
				}

				for field.Kind() == reflect.Ptr {
					field = field.Elem()
				}
			}

			if field.IsValid() && field.CanAddr() {
				setter(field, metaValue, context, record)
			}
		}
	}

	// Setup belongs_to / many_to_many Setter

	if meta.FieldStruct != nil {
		if relationship := meta.FieldStruct.Relationship; relationship != nil {
			if relationship.Kind == "belongs_to" || relationship.Kind == "many_to_many" {

				//commonSetter的参数后在 commonSetter返回的函数之内被调用
				meta.Setter = commonSetter(func(field reflect.Value, metaValue *MetaValue, context *core.Context, record interface{}) {
					var (
						scope         = context.GetDB().NewScope(record)
						indirectValue = reflect.Indirect(reflect.ValueOf(record))
					)
					primaryKeys := utils.ToArray(metaValue.Value)
					if metaValue.Value == nil {
						primaryKeys = []string{}
					}

					// associations not changed for belongs to
					if relationship.Kind == "belongs_to" && len(relationship.ForeignFieldNames) == 1 {
						oldPrimaryKeys := utils.ToArray(indirectValue.FieldByName(relationship.ForeignFieldNames[0]).Interface())
						// if not changed
						if fmt.Sprint(primaryKeys) == fmt.Sprint(oldPrimaryKeys) {
							return
						}

						// if removed
						if len(primaryKeys) == 0 {
							field := indirectValue.FieldByName(relationship.ForeignFieldNames[0])
							field.Set(reflect.Zero(field.Type()))
						}
					}

					// set current field value to blank
					field.Set(reflect.Zero(field.Type()))

					if len(primaryKeys) > 0 {
						// replace it with new value
						context.GetDB().Where(primaryKeys).Find(field.Addr().Interface())
					}

					// Replace many 2 many relations
					if relationship.Kind == "many_to_many" {
						if !scope.PrimaryKeyZero() {
							context.GetDB().Model(record).Association(meta.FieldName).Replace(field.Interface())
							field.Set(reflect.Zero(field.Type()))
						}
					}
				})
				return
			}
		}
	}

	field := reflect.Indirect(reflect.ValueOf(record)).FieldByName(fieldName)
	for field.Kind() == reflect.Ptr {
		if field.IsNil() {
			field.Set(utils.NewValue(field.Type().Elem()))
		}
		field = field.Elem()
	}

	if !field.IsValid() {
		return
	}

	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		meta.Setter = commonSetter(func(field reflect.Value, metaValue *MetaValue, context *core.Context, record interface{}) {
			field.SetInt(utils.ToInt(metaValue.Value))
		})
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		meta.Setter = commonSetter(func(field reflect.Value, metaValue *MetaValue, context *core.Context, record interface{}) {
			field.SetUint(utils.ToUint(metaValue.Value))
		})
	case reflect.Float32, reflect.Float64:
		meta.Setter = commonSetter(func(field reflect.Value, metaValue *MetaValue, context *core.Context, record interface{}) {
			field.SetFloat(utils.ToFloat(metaValue.Value))
		})
	case reflect.Bool:
		meta.Setter = commonSetter(func(field reflect.Value, metaValue *MetaValue, context *core.Context, record interface{}) {
			if utils.ToString(metaValue.Value) == "true" {
				field.SetBool(true)
			} else {
				field.SetBool(false)
			}
		})
	default:

		if _, ok := field.Addr().Interface().(sql.Scanner); ok {
			meta.Setter = commonSetter(func(field reflect.Value, metaValue *MetaValue, context *core.Context, record interface{}) {
				if scanner, ok := field.Addr().Interface().(sql.Scanner); ok {
					if metaValue.Value == nil && len(metaValue.MetaValues.Values) > 0 {
						decodeMetaValuesToField(meta.Resource, field, metaValue, context)
						return
					}

					if scanner.Scan(metaValue.Value) != nil {
						if err := scanner.Scan(utils.ToString(metaValue.Value)); err != nil {
							context.AddError(err)
							return
						}
					}
				}
			})
		} else if reflect.TypeOf("").ConvertibleTo(field.Type()) { //string

			meta.Setter = commonSetter(func(field reflect.Value, metaValue *MetaValue, context *core.Context, record interface{}) {
				field.Set(reflect.ValueOf(utils.ToString(metaValue.Value)).Convert(field.Type()))
			})
		} else if reflect.TypeOf([]string{}).ConvertibleTo(field.Type()) {
			meta.Setter = commonSetter(func(field reflect.Value, metaValue *MetaValue, context *core.Context, record interface{}) {
				field.Set(reflect.ValueOf(utils.ToArray(metaValue.Value)).Convert(field.Type()))
			})
		} else if _, ok := field.Addr().Interface().(*time.Time); ok {
			meta.Setter = commonSetter(func(field reflect.Value, metaValue *MetaValue, context *core.Context, record interface{}) {
				if str := utils.ToString(metaValue.Value); str != "" {
					if newTime, err := utils.ParseTime(str, context); err == nil {
						field.Set(reflect.ValueOf(newTime))
					}
				} else {
					field.Set(reflect.Zero(field.Type()))
				}
			})
		}
	}
}

//查找一个model struct下面的一个内嵌的struct， 返回这个内嵌的struct
func getNestedModel(value interface{}, fieldName string, context *core.Context) interface{} {
	model := reflect.Indirect(reflect.ValueOf(value))
	fields := strings.Split(fieldName, ".")
	for _, field := range fields[:len(fields)-1] {
		//当前model可否是struct,slice等可获取地址的
		if model.CanAddr() {
			submodel := model.FieldByName(field)
			if context != nil && context.GetDB() != nil && context.GetDB().NewRecord(submodel.Interface()) && !context.GetDB().NewRecord(model.Addr().Interface()) {
				if submodel.CanAddr() {
					context.GetDB().Model(model.Addr().Interface()).Association(field).Find(submodel.Addr().Interface())
					model = submodel
				} else {
					break
				}
			} else {
				model = submodel
			}
		}
	}

	if model.CanAddr() {
		return model.Addr().Interface()
	}
	return nil
}
