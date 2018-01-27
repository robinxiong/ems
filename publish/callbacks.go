package publish

import "github.com/jinzhu/gorm"

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

			//开始draft模式
			if ensureDraftMode{
				scope.Set("publish:force_draft_table", true)
				//指定一个_dratf table, scope的相关操作都是针对这个数所库表
				scope.Search.Table(DraftTableName(scope.TableName()))
				// 当从draft talbes中更新数据时，用来设置publish status状态
				//db中是否设置了publish:draft_mode

				if IsDraftMode(scope.DB()) {
					//db中是否设置了publish:publish_event变量，如果设置了, 则创建publish:creating_publish_event,并设置true,
					//设置PublishStatus列为true
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
