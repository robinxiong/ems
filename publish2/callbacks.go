package publish2

import (
	"github.com/jinzhu/gorm"
	"reflect"
	"ems/core/utils"
	"time"
	"fmt"
)

const (
	ModeOff             = "off"
	ModeReverse         = "reverse"
	VersionMode         = "publish:version:mode"
	VersionNameMode     = "publish:version:name"
	VersionMultipleMode = "multiple"

	ScheduleMode     = "publish:schedule:mode"
	ComingOnlineMode = "coming_online"
	GoingOfflineMode = "going_offline"
	ScheduledTime    = "publish:schedule:current"
	ScheduledStart   = "publish:schedule:start"
	ScheduledEnd     = "publish:schedule:end"

	VisibleMode = "publish:visible:mode"
)



func RegisterCallbacks(db *gorm.DB) {
	//todo: querycallback
	//如果没有设置version_name则使用默认的值，同时更新priority列,  如果是ShareableVersionModel,则更新它field为IsBlank=false
	db.Callback().Create().Before("gorm:begin_transaction").Register("publish:versions", createCallback)
	//只更新priority列
	db.Callback().Update().Before("gorm:begin_transaction").Register("publish:versions", updateCallback)
	//检查db中是否设置了publish:version:name，获取versionName, 然后删除指定的version_name的记录
	//同时在默认的deleteCallback中，会把主键信息也会写入到where条件scope.CombinedConditionSql()
	db.Callback().Delete().Before("gorm:begin_transaction").Register("publish:versions", deleteCallback)
}


func createCallback(scope *gorm.Scope){
	//查看是否实现了version接口，即model中是否可以设置和获取version
	if IsVersionableModel(scope.Value) {

		//如果实现了，即包含了publish2.Version, 查看数据库表是否包含version_name
		if field, ok := scope.FieldByName("VersionName"); ok {
			//如果这个字段的值为空，设置一个默认的值
			if field.IsBlank {
				field.Set(DefaultVersionName)
			}
		}
		//接着更新version的priority, 优先级
		//是否包含schedule字段ScheduledStartAt, 如果设置了开始发布的时间，则设置version_priority的值为2006-01-02T15:04:05Z07:00_versionName
		updateVersionPriority(scope)
	}

	//是否实现了version中的ShareableVersionInterface
	//如果model中的ShareableVersion的值为空，也设置它为空
	if IsShareableVersionModel(scope.Value){
		if field, ok := scope.FieldByName("VersionName"); ok {
			field.IsBlank = false
		}
	}
}

func updateCallback(scope *gorm.Scope) {
	if IsVersionableModel(scope.Value) {
		updateVersionPriority(scope)
	}
}

func deleteCallback(scope *gorm.Scope) {
	if versionName, ok := scope.DB().Get(VersionNameMode); ok && versionName != "" {
		if IsVersionableModel(scope.Value) || IsShareableVersionModel(scope.Value) {
			scope.Search.Where("version_name = ?", versionName)
		}
	}
}

//是否实现了version接口，即get和set version
func IsVersionableModel(model interface{}) (ok bool){
	if model != nil {
		_, ok = reflect.New(utils.ModelType(model)).Interface().(VersionableInterface)
	}
	return
}

func IsShareableVersionModel(model interface{}) (ok bool) {
	if model != nil {
		_, ok = reflect.New(utils.ModelType(model)).Interface().(ShareableVersionInterface)
	}
	return
}

func updateVersionPriority(scope *gorm.Scope) {
	if field, ok := scope.FieldByName("VersionPriority"); ok {
		var scheduledTime *time.Time //开始时间
		var versionName string

		//是否包含schedule字段ScheduledStartAt, ScheduledEndAt
		//如果没有，南使用当前的时间
		if scheduled, ok := scope.Value.(ScheduledInterface); ok {
			scheduledTime = scheduled.GetScheduledStartAt()
		}

		if scheduledTime == nil {
			unix := time.Unix(0, 0) //1970-01-01 08:00:00 +0800 C
			scheduledTime = &unix
		}

		if versionable, ok := scope.Value.(VersionableInterface); ok {
			versionName = versionable.GetVersionName()
		}

		//2006-01-02T15:04:05Z07:00_versionName
		//每一个版本的发布时间
		priority := fmt.Sprintf("%v_%v", scheduledTime.UTC().Format(time.RFC3339), versionName)
		field.Set(priority)

	}
}

//查询会使用到以下方法