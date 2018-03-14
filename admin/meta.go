package admin

import (
	"database/sql"
	"ems/core"
	"ems/core/resource"
	"ems/core/utils"
	"ems/roles"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"time"
)

type MetaConfigInterface interface {
	resource.MetaConfigInterface
}

type Meta struct {
	*resource.Meta
	Name            string
	FieldName       string
	Type            string //hidden_primary_key meta admin/resource.go GetMeta
	Setter          func(record interface{}, metaValue *resource.MetaValue, context *core.Context)
	Valuer          func(record interface{}, context *core.Context) (result interface{})
	FormattedValuer func(record interface{}, context *core.Context) (result interface{})
	Permission      *roles.Permission
	Config          MetaConfigInterface
	//表单的label
	Label      string //默认为meta.Name的可读模式 e.g.Item Name
	Collection interface{}

	baseResource *Resource //所对应的资源
	Resource     *Resource
}

func (meta *Meta) configure() {

	if meta.Meta == nil {
		meta.Meta = &resource.Meta{
			Name:            meta.Name,
			FieldName:       meta.FieldName,
			Setter:          meta.Setter,
			Valuer:          meta.Valuer,
			FormattedValuer: meta.FormattedValuer,
			BaseResource:    meta.baseResource,
			Resource:        meta.Resource,
			Permission:      meta.Permission,
			Config:          meta.Config,
		}
	} else {
		meta.Meta.Name = meta.Name
		meta.Meta.FieldName = meta.FieldName
		meta.Meta.Setter = meta.Setter
		meta.Meta.Valuer = meta.Valuer
		meta.Meta.FormattedValuer = meta.FormattedValuer
		meta.Meta.BaseResource = meta.baseResource
		meta.Meta.Resource = meta.Resource
		meta.Meta.Permission = meta.Permission
		meta.Meta.Config = meta.Config
	}
	//core/resource/meta.go
	//设置FildName和FieldStruct
	meta.PreInitialize()
	if meta.FieldStruct != nil {
		//meta_test.go中，User聚合了media.OSS，而media.Base是实现了这个接口
		//所以先调用Base的ConfigureMetaBeforeInitialize
		if injector, ok := reflect.New(meta.FieldStruct.Struct.Type).Interface().(resource.ConfigureMetaBeforeInitializeInterface); ok {
			injector.ConfigureMetaBeforeInitialize(meta)
		}
	}
	//设置Valuer, Setter
	meta.Initialize()

	if meta.Label == "" {
		meta.Label = utils.HumanizeString(meta.Name)
	}

	//设置 meta 的类型 type
	var fieldType reflect.Type              //获取meta所对应字段类型
	var hasColumn = meta.FieldStruct != nil //在PreInitialize中设置

	if hasColumn {
		fieldType = meta.FieldStruct.Struct.Type
		for fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}
	}

	if hasColumn {
		//只在meta.Type为空的时候调用, 比如之前在getMeta获取到PrimaryField时，设置为hidden_primary_key, 所以不用设置
		if meta.Type == "" {
			//创建一个字段类型,
			if _, ok := reflect.New(fieldType).Interface().(sql.Scanner); ok {

				//字段是一个struct类型，同时实现了sql.Scanner，即有Scan方法
				if fieldType.Kind() == reflect.Struct {
					fieldType = reflect.Indirect(reflect.New(fieldType)).Field(0).Type()
				}
			}
			//slice或者指定了两列Profile,ProfileID等, meta_test.go TestRelationFieldMetaType
			if relationship := meta.FieldStruct.Relationship; relationship != nil {

				if relationship.Kind == "has_one" {
					meta.Type = "single_edit"
				} else if relationship.Kind == "has_many" {
					meta.Type = "collection_edit"
				} else if relationship.Kind == "belongs_to" {
					meta.Type = "select_one"
				} else if relationship.Kind == "many_to_many" {
					meta.Type = "select_many"
				}
			} else {
				//常规的类型，比如string, int, bool
				switch fieldType.Kind() {
				case reflect.String:
					var tags = meta.FieldStruct.TagSettings
					if size, ok := tags["SIZE"]; ok {
						//如果数据库表的字段大于255, 则为textarea
						if i, _ := strconv.Atoi(size); i > 255 {
							meta.Type = "text"
						} else {
							meta.Type = "string"
						}

					} else if text, ok := tags["TYPE"]; ok && text == "text" { //在tag中指定了TYPE
						meta.Type = "text"
					} else {
						meta.Type = "string"
					}
				case reflect.Bool:
					meta.Type = "checkbox"
				default:
					if regexp.MustCompile(`^(.*)?(u)?(int)(\d+)?`).MatchString(fieldType.Kind().String()) {
						meta.Type = "number"
					} else if regexp.MustCompile(`^(.*)?(float)(\d+)?`).MatchString(fieldType.Kind().String()) {
						meta.Type = "float"
					} else if _, ok := reflect.New(fieldType).Interface().(*time.Time); ok {
						meta.Type = "datetime"
					} else {
						//struct
						if fieldType.Kind() == reflect.Struct {
							meta.Type = "single_edit"
							log.Println(meta.Type)
						} else if fieldType.Kind() == reflect.Slice {
							refelectType := fieldType.Elem()
							for refelectType.Kind() == reflect.Ptr {
								refelectType = refelectType.Elem()
							}
							if refelectType.Kind() == reflect.Struct {
								meta.Type = "collection_edit"
							}
						}
					}
				}
			}

		}

	}

	// Set meta Resource
	{
		if hasColumn {

			if meta.Resource == nil {
				var result interface{}

				if fieldType.Kind() == reflect.Struct {
					result = reflect.New(fieldType).Interface()
				} else if fieldType.Kind() == reflect.Slice {
					refelectType := fieldType.Elem()
					for refelectType.Kind() == reflect.Ptr {
						refelectType = refelectType.Elem()
					}
					if refelectType.Kind() == reflect.Struct {
						result = reflect.New(refelectType).Interface()
					}
				}

				if result != nil {
					res := meta.baseResource.NewResource(result)
					meta.Resource = res
					//将subResource的权限赋值给meta
					meta.Meta.Permission = meta.Meta.Permission.Concat(res.Config.Permission)
				}

			}

			if meta.Resource != nil {
				permission := meta.Resource.Permission.Concat(meta.Meta.Permission)
				meta.Meta.Resource = meta.Resource
				meta.Resource.Permission = permission
				meta.SetPermission(permission)
			}
		}
	}

	meta.FieldName = meta.GetFieldName()

	// call meta config's ConfigureMetaInterface
	if meta.Config != nil {
		meta.Config.ConfigureMeta(meta)
	}

}

// HasPermission check has permission or not
func (meta Meta) HasPermission(mode roles.PermissionMode, context *core.Context) bool {
	var roles = []interface{}{}
	for _, role := range context.Roles {
		roles = append(roles, role)
	}
	if meta.Permission != nil {
		return meta.Permission.HasPermission(mode, roles...)
	}

	if meta.Resource != nil {
		return meta.Resource.HasPermission(mode, context)
	}

	if meta.baseResource != nil {
		return meta.baseResource.HasPermission(mode, context)
	}

	return true
}

// DBName get meta's db name, used in index page for sorting, example: views/index/table.tmpl
func (meta *Meta) DBName() string {
	if meta.FieldStruct != nil {
		return meta.FieldStruct.DBName
	}
	return ""
}
