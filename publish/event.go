package publish

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"fmt"
)

// EventInterface defined methods needs for a publish event
type EventInterface interface {
	Publish(db *gorm.DB, event PublishEventInterface) error
	Discard(db *gorm.DB, event PublishEventInterface) error
}

//注册所有发布事件, 保存到map中
var events = map[string]EventInterface{}


// RegisterEvent register publish event
func RegisterEvent(name string, event EventInterface) {
	events[name] = event
}

// PublishEvent default publish event model, 主要是更新publish_event表, 同时调用EventInterface的publish方法
// 实现PublishEventInterface, 参见publish.go
// 它是一个事件的基类，它的子类还需要实现EventInterface接口
type PublishEvent struct {
	gorm.Model
	Name          string
	Description   string
	Argument      string `sql:"size:65532"`
	PublishStatus bool
	PublishedBy   string
}

//发布的时假需要获取是谁发布的
func getCurrentUser(db *gorm.DB) (string, bool) {
	if user, hasUser := db.Get("ems:current_user"); hasUser {
		var currentUser string
		//如果包含主键，则返回主键的值,
		//否则，返回当前用户对像的字符串形式
		if primaryField := db.NewScope(user).PrimaryField(); primaryField != nil {
			currentUser = fmt.Sprintf("%v", primaryField.Field.Interface())
		} else {
			currentUser = fmt.Sprintf("%v", user)
		}

		return currentUser, true
	}
	return "", false
}


// Publish 发布数据
func (publishEvent *PublishEvent) Publish(db *gorm.DB) error {
	//首先检查是否在events中注册
	if event, ok := events[publishEvent.Name]; ok {
		err := event.Publish(db, publishEvent)
		if err == nil {
			//更新publish_status为PUBLISHED (false)
			var updateAttrs = map[string]interface{}{"PublishStatus": PUBLISHED}
			if user, hasUser := getCurrentUser(db); hasUser {
				updateAttrs["PublishedBy"] = user
			}
			//更新publishEvent表
			err = db.Model(publishEvent).Update(updateAttrs).Error
		}
	}
	return errors.New("event not found")
}

func (publishEvent *PublishEvent) Discard(db *gorm.DB) error {
	if event, ok := events[publishEvent.Name]; ok {
		err := event.Discard(db, publishEvent)
		if err == nil {
			var updateAttrs = map[string]interface{}{"PublishStatus": PUBLISHED}
			if user, hasUser := getCurrentUser(db); hasUser {
				updateAttrs["PublishedBy"] = user
			}
			err = db.Model(publishEvent).Update(updateAttrs).Error
		}
		return err
	}
	return errors.New("event not found")
}


