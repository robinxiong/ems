package admin

import (
	"ems/core"
	"ems/core/resource"
	"ems/core/utils"
	"ems/roles"

	"fmt"
	"net/http"
	"reflect"

	"github.com/jinzhu/gorm"
	"github.com/jinzhu/inflection"
)

// Config resource config struct
type Config struct {
	Name       string   //指定Resource的名字字，默认为model struct的名字OrderItem，则为Order Item
	Menu       []string //父菜单的路径 "中国" "广东"
	Permission *roles.Permission
	Themes     []ThemeInterface
	Priority   int  //菜单中的优先级
	Singleton  bool //是否是当个Resource, 如果不是，则当前资源显示的菜单名称将为复数，如Colors, 同时在route.go中，注册resource路由时有用
	Invisible  bool //是否显示资源,默认值为false, 即显示资源, 则会将当前Resource添加到菜单
	PageCount  int
}

// Resource 在admin中是一个非常重要的概念，所以的数据库表model都抽像为一个资源, admin生成的管理接口，都是基于它的
type Resource struct {
	*resource.Resource //core/resource 用于将form数据写入到数据库
	Config             *Config
	ParentResource     *Resource //可以参考route, RegisterRoute中，将使用到ParentResource
	SearchHandler      func(keyword string, context *core.Context) *gorm.DB
	params             string
	admin              *Admin    //每一个Resource保留一份admin对像
	mounted            bool      //将resource到路由时使用, newRouteHandler
	scopes             []*Scope  //当request中包含了scopes，则创建scopes来搜索数据库，跟filter的区别是，filter会包含要过滤的字段，而scope是直接定义好了，只需要提供值就可以
	filters            []*Filter //request中包含了filters["color"]="red" filters["size"]=1, 用于resource过滤搜索结果

	metas []*Meta //admin.AddResource和NewResource都会调用res.configure()，将数据库表中的字段，封装为Meta

	//用来组织resource的表格，比如product的基本信息有哪些字段, Organization包含Category, Gender, Seo Meta section
	sections struct {
		IndexSections                  []*Section //显示在index页面的section, 如果没有指定，将生成一个默认的Section
		OverriddingIndexAttrs          bool
		OverriddingIndexAttrsCallbacks []func()

		SortableAttrs                  *[]string //views/index/table.tmpl 是否可以按标题进行排序
	}
}

func (res *Resource) GetAdmin() *Admin {
	return res.admin
}



// NewResource initialize a new resource, won't add it to admin, just initialize it
// 创建一个新的Resource, 但不添加到admin.resources中,只是初始化， example: admin/meta.config 设置meta.Resource
func (res *Resource) NewResource(value interface{}, config ...*Config) *Resource {
	subRes := res.GetAdmin().newResource(value, config...)
	subRes.ParentResource = res
	subRes.configure()
	return subRes
}



// ToParam 将一个resource的名称,如Order Item转换为url param形式order_item，即替换空格为下划线
// 首先确定res.params是不为空，如果有res.params则直接返回， 否则调用model struct的ToParam(model struct实现了ToParam方法)
// 如果没有实现ToParam方法，则将res.Name转换成url允许的形式
// 用法可以参考AddResource中添加AddMenu (即生成菜单的RelativePath
// 第二个用法是route中RegisterRoute调用, 即向router中注册路径
func (res *Resource) ToParam() string {
	if res.params == "" {
		if value, ok := res.Value.(interface {
			ToParam() string
		}); ok {
			res.params = value.ToParam()
		} else {
			if res.Config.Singleton == true {
				res.params = utils.ToParamString(res.Name)
			} else {
				res.params = utils.ToParamString(inflection.Plural(res.Name))
			}
		}
	}
	return res.params
}

// ParamIDName ParamIDName return param name for primary key like :product_id
// 用法RegisterRoute
func (res Resource) ParamIDName() string {
	return fmt.Sprintf(":%v_id", inflection.Singular(utils.ToParamString(res.Name)))
}

// GetPrimaryValue 从Request中获取:color_id的名称, color为resource的名称
func (res Resource) GetPrimaryValue(request *http.Request) string {
	if request != nil {
		return request.URL.Query().Get(res.ParamIDName())
	}
	return ""
}

//IndexAttrs 用来设置显示在index page的属性, func_map的indexSections方法，传递的参数为空, 则调用setSections时，生成一个默认的section
//否则，则使用EditAttrs设置的Section
func (res *Resource) IndexAttrs(values ...interface{}) []*Section {
	//todo:overrideIndexAttr
	res.setSections(&res.sections.IndexSections, values...)
	return res.sections.IndexSections
}

//将数据库字段，封装成Meta
func (res *Resource) configure() {
	//首先确认内嵌的字段，是否实现了ConfigureResourceInterface接口，如果没有，则到子字段
	var configureModel func(value interface{})

	configureModel = func(value interface{}) {
		modelType := utils.ModelType(value)
		for i := 0; i < modelType.NumField(); i++ {
			if fieldStruct := modelType.Field(i); fieldStruct.Anonymous {
				if injector, ok := reflect.New(fieldStruct.Type).Interface().(resource.ConfigureResourceInterface); ok {
					injector.ConfigureResource(res)
				} else {
					configureModel(reflect.New(fieldStruct.Type).Interface())
				}
			}
		}
	}

	configureModel(res.Value)

	//遍历model struct中的每一个字段, 如果它的类型为struct, 同时实现了ConfigureMetaInterface或者ConfigureMetaBeforeInitializeInterface, 封装为Meta
	scope := gorm.Scope{Value: res.Value}

	for _, field := range scope.Fields() {
		//field.StructField是qorm用来描述field的一个struct类型，它包含一个字段的db名称，是否为主键等
		//StructField.Struct是表示它在model struct中的定义, 比如类型，名称，是否为内嵌的, tag等
		//如果这个字段在model struct的类型为一个struct
		if field.StructField.Struct.Type.Kind() == reflect.Struct {
			fieldData := reflect.New(field.StructField.Struct.Type).Interface() //以interface{}返回
			//是否实现了resource的两个接口
			_, configureMetaBeforeInitialize := fieldData.(resource.ConfigureMetaBeforeInitializeInterface)
			_, configureMeta := fieldData.(resource.ConfigureMetaInterface)

			//如果实现了作何一个，则对它进行封装
			if configureMetaBeforeInitialize || configureMeta {
				res.Meta(&Meta{Name: field.Name})
			}
		}
	}

	if injector, ok := res.Value.(resource.ConfigureResourceInterface); ok {
		injector.ConfigureResource(res)
	}

}

func (res *Resource) Meta(meta *Meta) *Meta {
	if oldMeta := res.GetMeta(meta.Name); oldMeta != nil {
		if meta.Type != "" {
			oldMeta.Type = meta.Type
			oldMeta.Config = nil
		}

		if meta.Label != "" {
			oldMeta.Label = meta.Label
		}

		if meta.FieldName != "" {
			oldMeta.FieldName = meta.FieldName
		}

		if meta.Setter != nil {
			oldMeta.Setter = meta.Setter
		}

		if meta.Valuer != nil {
			oldMeta.Valuer = meta.Valuer
		}

		if meta.FormattedValuer != nil {
			oldMeta.FormattedValuer = meta.FormattedValuer
		}

		if meta.Resource != nil {
			oldMeta.Resource = meta.Resource
		}

		if meta.Permission != nil {
			oldMeta.Permission = meta.Permission
		}

		if meta.Config != nil {
			oldMeta.Config = meta.Config
		}

		if meta.Collection != nil {
			oldMeta.Collection = meta.Collection
		}
		meta = oldMeta
	} else {
		res.metas = append(res.metas, meta)
		meta.baseResource = res
	}

	meta.configure()

	return meta
}

//返回resource中model struct所对应的数据库字段
func (res *Resource) allAttrs() []string {
	var attrs []string
	scope := &gorm.Scope{Value: res.Value}

Fields:
	for _, field := range scope.GetModelStruct().StructFields {

		//当前字段是否为Meta对像, 则跳过
		for _, meta := range res.metas {
			if field.Name == meta.FieldName {
				attrs = append(attrs, meta.Name)
				continue Fields
			}
		}

		//如果field是外键，则跳过
		if field.IsForeignKey {
			continue
		}

		//跳过CreatedAt, UpdatedAt, DeletedAt等字段的设置
		for _, value := range []string{"CreatedAt", "UpdatedAt", "DeletedAt"} {
			if value == field.Name {
				continue Fields
			}
		}

		//将正常的字段添加到attrs中
		if (field.IsNormal || field.Relationship != nil) && !field.IsIgnored {
			attrs = append(attrs, field.Name)
			continue
		}

		//model struct字段的类型为指针，或者数组，或者struct, 也添加到attrs中
		fieldType := field.Struct.Type
		for fieldType.Kind() == reflect.Ptr || fieldType.Kind() == reflect.Slice {
			fieldType = fieldType.Elem()
		}

		if fieldType.Kind() == reflect.Struct {
			attrs = append(attrs, field.Name)
		}
	}
MetaIncluded:
	//遍历metas, 如果attr等于meta, 则跳过，到下一个metas, 否则将meta添加到attrs中
	for _, meta := range res.metas {
		for _, attr := range attrs {
			if attr == meta.FieldName || attr == meta.Name {
				continue MetaIncluded
			}
		}
		attrs = append(attrs, meta.Name)
	}

	return attrs
}
//allowedSections
//如果name没有存在于res.Metas中，则从model struct中查找到字段，在封装成Meta, 在调用meta.configure进行配置
func (res *Resource) GetMeta(name string) *Meta{
	var fallbackMeta *Meta
	//是否有在resource中定义Meta, 即调用res.Meta()方法，添加
	//或者model struct的字段实现了configureResourceInterface接口
	for _, meta := range res.metas {
		if meta.Name == name {
			return meta
		}

		if meta.GetFieldName() == name {
			fallbackMeta = meta
		}
	}

	if fallbackMeta == nil {
		if field, ok := res.GetAdmin().DB.NewScope(res.Value).FieldByName(name); ok {
			meta := &Meta{Name: name, baseResource: res}
			if field.IsPrimaryKey {
				meta.Type = "hidden_primary_key"
			}
			meta.configure()
			res.metas = append(res.metas, meta)
			return meta
		}
	}

	return fallbackMeta

}

//确认当前用户哪些section可以访问
func (res *Resource) allowedSections(sections []*Section, context *Context, roles ...roles.PermissionMode) []*Section {
	var newSections []*Section
	for _, section := range sections {
		newSection := Section{Resource: section.Resource, Title: section.Title}
		var editableRows [][]string
		for _, row := range section.Rows {
			var editableColumns []string
			for _, column := range row {
				for _, role := range roles {
					meta := res.GetMeta(column)
					if meta != nil && meta.HasPermission(role, context.Context) {
						editableColumns = append(editableColumns, column)
						break
					}
				}
			}
			if len(editableColumns) > 0 {
				editableRows = append(editableRows, editableColumns)
			}
		}
		newSection.Rows = editableRows
		newSections = append(newSections, &newSection)
	}
	return newSections
}

