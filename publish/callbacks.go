package publish

import "github.com/jinzhu/gorm"


//isProductionModeAndNewScope 用来检测当前db是publish模式(或者没有设置)，同时table实现了publishInterface
//如果满足条件，则返回true, 则一个新的scope
func isProductionModeAndNewScope(scope *gorm.Scope) (isProduction bool, clone *gorm.Scope){
	if !IsDraftMode(scope.DB()){
		//在setTableAndPublishStatus中设置publish:supported_model
		if _, ok := scope.InstanceGet("publish:supported_model"); ok {
			table := OriginalTableName(scope.TableName())
			clone := scope.New(scope.Value)
			clone.Search.Table(table)
			return true, clone
		}
	}
	return false, nil
}


//setTableAndPublishStatus 当db为draft模式时， 用于更新publish_status列为DIRTY(true), 如果db中还设置了publish:publish_event变量， 则设置scope的publish:creating_publish_event变量为true, 而不更新publish_status列
// ensureDraftMode用来如果为true, 则强制先写入_draft数据表, 然后在提交事务时，如果是publish模式，则在复制到publish表中， 则如果是draft模式，则不在复制到publish表中
func setTableAndPublishStatus(ensureDraftMode bool) func(scope *gorm.Scope) {
	return func (scope *gorm.Scope) {
		if scope.Value == nil {
			return
		}
		//要保存的model是否实现了publishInterface接口
		if IsPublishableModel(scope.Value){
			//InstanceSet只针对当前scope有效，instanceID为scope地址，以及scope.db地址，对于有的回调，比如saving associations callback无效，它新建了一个db以及scope
			//设置publish:supported_model,用于检测当前scope是否支持publish
			scope.InstanceSet("publish:supported_model", true)

			//开启draft模式, 查看publish.New方法，即不管db是否为draft模式，都先写入_draft后缀， 比如products_draft
			if ensureDraftMode{
				scope.Set("publish:force_draft_table", true)
				//指定一个_dratf table, scope的相关操作都是针对这个数所库表
				scope.Search.Table(DraftTableName(scope.TableName()))
				// 当从draft talbes中更新数据时，用来设置publish status状态
				//db中是否设置了publish:draft_mode

				if IsDraftMode(scope.DB()) {
					//db中是否设置了publish:publish_event变量，如果设置了, 则创建publish:creating_publish_event,并设置true,
					//设置PublishStatus列为true
					//查看 publish_event_test.go TestCreateNewEvent
					if _, ok := scope.DB().Get(publishEvent); ok {
						scope.InstanceSet("publish:creating_publish_event", true)
					} else {
						//gorm:update_attrs会在db.updateColumns中设置，用于指定更新哪些列
						if attrs, ok := scope.InstanceGet("gorm:update_attrs"); ok {
							updateAttrs := attrs.(map[string]interface{})
							updateAttrs["publish_status"] = DIRTY
							scope.InstanceSet("gorm:update_attrs", updateAttrs)
						} else {
							scope.SetColumn("PublishStatus", DIRTY)
						}
					}
				}
			}
		}
	}
}
//setTableAndPublishStatus回调中设置了scope的tableName为_draft(不管db是否为publish或者draft模式)
//如果是production,即db中没有设置draft模式，则复制到production表中创建一个相同的记录
func syncCreateFromProductionToDraft(scope *gorm.Scope){
	if !scope.HasError() {
		//克隆一个scope, 重新执行createCallback(gorm/callback_create.go), 即数据库表中创建一个记录
		if ok, clone := isProductionModeAndNewScope(scope); ok{
			scope.DB().Callback().Create().Get("gorm:create")(clone)
		}
	}
}

//如果db中设置了publish:publish_event变量, 则不更新数据库表publish_status列为DIRTY
//而是将它保存到publish_events表中
func createPublishEvent(scope *gorm.Scope){
	if _, ok := scope.InstanceGet("publish:creating_publish_event"); ok {
		//ublish_event_test.go TestCreateProductWithPublishAllEvent, 而如果不是event对像，则忽略
		if event, ok := scope.Get(publishEvent); ok {
			if event, ok := event.(*PublishEvent); ok {
				event.PublishStatus = DIRTY
				scope.Err(scope.NewDB().Save(&event).Error)
			}
		}
	}
}
