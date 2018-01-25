package publish

import (
	"github.com/jinzhu/gorm"
	"ems/core/utils"
	"reflect"
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
func New(db *gorm.DB) *Publish{
	tableHandler := gorm.DefaultTableNameHandler //默认直接返回tableName
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultName string) string {
		tableName := tableHandler(db, defaultName)
		//自定义model struct对应的表名字
		if db != nil {
			//db.Value为设置了model的值
			if IsPublishableModel(db.Value) {
				typ := utils.ModelType(db.Value)
				//如果没有缓存此model, 因injectedJoinTableHandler为bool值，如果没有找到，则为零值false
				if !injectedJoinTableHandler[typ] {
					injectedJoinTableHandler[typ] = true //设置缓存
					scope := db.NewScope(db.Value)
					for _, field := range scope.GetModelStruct().StructFields {
						if many2many := utils.ParseTagOption(field.Tag.Get("gorm"))["MANY2MANY"]; many2many!= nil {
							db.Set
						}
					}

				}
			}
		}
		return tableName
	}
	return nil
}


// IsPublishableModel check if current model is a publishable
// 如果一个struct包含了Status, 则实现了publishInterface
func IsPublishableModel(model interface{}) (ok bool) {
	if model != nil {
		_, ok = reflect.New(utils.ModelType(model)).Interface().(publishInterface)
	}
	return
}
