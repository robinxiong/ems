package publish

import (
	"ems/core/utils"
	"reflect"
	"strings"

	"github.com/jinzhu/gorm"
)

const (
	//publish status
	PUBLISHED = false
	//publish status
	DIRTY = true
	//设置db模式
	publishDraftMode = "publish:draft_mode"
	publishEvent     = "publish:publish_event"
)

type publishInterface interface {
	GetPublishStatus() bool
	SetPublishStatus(bool)
}

// PublishEventInterface defined publish event itself's interface
type PublishEventInterface interface {
	Publish(*gorm.DB) error
	Discard(*gorm.DB) error
}

//Status publish_status 实现了publishInterface
type Status struct {
	PublishStatus bool
}

// GetPublishStatus get publish status
func (s Status) GetPublishStatus() bool {
	return s.PublishStatus
}

// SetPublishStatus set publish status
func (s *Status) SetPublishStatus(status bool) {
	s.PublishStatus = status
}

type Publish struct {
	DB *gorm.DB
}

//缓存生成的model表名
var injectedJoinTableHandler = map[reflect.Type]bool{}

//初妈化一个Publish instance

func New(db *gorm.DB) *Publish {
	/*
			我们知道，db.AutoMigrate(&Product{}), 它会调用到model_struct.go的TableName方法，这个方法会在最后的位置调用
		DefaultTableNameHandler, DefaultTableNameHandler通常直接返回传入的表名，不做作何处理，而在这里，我们需要对publish类型的DB
		调用AutoMigrate返回 _draft等后缀, 就需要修改这个函数
	*/
	tableHandler := gorm.DefaultTableNameHandler //默认直接返回tableName
	//修改gorm包的默认返回表的函数，只要调用了pulish.New(), 其它包都将使用这个函数来返回表名
	//当删除表时，也会生成many2many都关联表

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultName string) string {

		tableName := tableHandler(db, defaultName) //调用gorm中定义的DefaultTableNameHandler
		//自定义model struct对应的表名字
		if db != nil {
			//db.Value为设置了model的值实现了publishInterface
			if IsPublishableModel(db.Value) {
				typ := utils.ModelType(db.Value)
				//如果没有缓存此model, 因injectedJoinTableHandler为bool值，如果没有找到，则为零值false
				if !injectedJoinTableHandler[typ] {
					injectedJoinTableHandler[typ] = true //设置缓存
					scope := db.NewScope(db.Value)

					for _, field := range scope.GetModelStruct().StructFields {
						//如果包含了many2many的字段，我们需要先创建关联表
						if many2many := utils.ParseTagOption(field.Tag.Get("gorm"))["MANY2MANY"]; many2many != "" {
							//调用SetJoinTableHandler, 找到此例，比如此例的名称为 Categories []Category `gorm:"many2many:product_categories"`
							//source为 Product
							//destination为Category
							//调用publishJoinTableHandler.Setup, 它的第一个参数为field.Relationship， relationship是通过scope.GetModelStruct()获得
							//第二个参数many2many, product_categories, 第三个参数为source, 第四个为destination
							//如果db设置了public_draft, 则publishJoinTableHandler返回product_categories_draft
							//但没有必要在SetJoinTableHandler方法中调用s.Table(table).AutoMigrate(handler)
							db.SetJoinTableHandler(db.Value, field.Name, &publishJoinTableHandler{})
						}
					}
					//创建整个表，对于dropTable如果它存在many2many它也会创建子表, 比如drop(&Product{}) 它会先创建Product, product_categories product_languages, 然后在删除product


				}

				var forceDraftTable bool
				if forceDraft, ok := db.Get("publish:force_draft_table"); ok {
					if forceMode, ok := forceDraft.(bool); ok && forceMode {
						forceDraftTable = true
					}
				}

				if IsDraftMode(db) || forceDraftTable {
					return DraftTableName(tableName)
				}
			}
		}
		return tableName
	}

	//创建PublishEvent表，记录publish 和 discard 事件
	db.AutoMigrate(&PublishEvent{})

	//注册publish回调函数
	//在事件开启前，设置db为
	db.Callback().Create().Before("gorm:begin_transaction").Register("publish:set_table_to_draft", setTableAndPublishStatus(true))
	//在提交事务前，注册回调
	db.Callback().Create().Before("gorm:commit_or_rollback_transaction").Register("publish:sync_to_production_after_create", syncCreateFromProductionToDraft)
	//提交事件前，注册回调
	db.Callback().Create().Before("gorm:commit_or_rollback_transaction").Register("gorm:create_publish_event", createPublishEvent)
	return &Publish{DB: db}
}

// IsDraftMode 检查db是否设置了publish:draft_mode
func IsDraftMode(db *gorm.DB) bool {
	if draftMode, ok := db.Get(publishDraftMode); ok {
		if isDraft, ok := draftMode.(bool); ok && isDraft {
			return true
		}
	}
	return false
}

// IsPublishableModel check if current model is a publishable
// 如果一个struct包含了Status, 则实现了publishInterface
func IsPublishableModel(model interface{}) (ok bool) {
	if model != nil {
		_, ok = reflect.New(utils.ModelType(model)).Interface().(publishInterface)
	}
	return
}

func DraftTableName(table string) string {
	return OriginalTableName(table) + "_draft"
}

// OriginalTableName get original table name of passed in string
func OriginalTableName(table string) string {
	return strings.TrimSuffix(table, "_draft")
}

// ProductionDB get db in production mode
func (pb Publish) ProductionDB() *gorm.DB {
	return pb.DB.Set(publishDraftMode, false)
}

// DraftDB get db in draft mode
func (pb Publish) DraftDB() *gorm.DB {
	return pb.DB.Set(publishDraftMode, true)
}

// AutoMigrate run auto migrate in draft tables
func (pb *Publish) AutoMigrate(values ...interface{}) {
	for _, value := range values {
		tableName := pb.DB.NewScope(value).TableName()
		//Table方法用于返回一个db, 设置它的db.Search.table = 源表名_draft
		//创建这个数据库表
		pb.DraftDB().Table(DraftTableName(tableName)).AutoMigrate(value)
	}
}
