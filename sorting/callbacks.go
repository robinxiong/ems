package sorting

import (
	"ems/l10n"
	"ems/publish"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/jinzhu/gorm"
	"log"
)

func initalizePosition(scope *gorm.Scope) {

	if !scope.HasError() {
		if _, ok := scope.Value.(sortingInterface); ok {
			var lastPosition int
			//找到当前position最高的值，然后设置新行的Position+1
			scope.NewDB().Set("l10n:mode", "locale").Model(modelValue(scope.Value)).Select("position").Order("position DESC").Limit(1).Row().Scan(&lastPosition)
			scope.SetColumn("Position", lastPosition+1)
		}
	}
}

//获取value的值，并且以Interface返回
func modelValue(value interface{}) interface{} {
	reflectValue := reflect.Indirect(reflect.ValueOf(value))
	if reflectValue.IsValid() {
		typ := reflectValue.Type()

		if reflectValue.Kind() == reflect.Slice {
			typ = reflectValue.Type().Elem()
			if typ.Kind() == reflect.Ptr {
				typ = typ.Elem()
			}
		}

		return reflect.New(typ).Interface()
	}
	return nil
}

//对于draft model来说，需要在publish_event中添加一条reorder的记录
func createPublishEvent(db *gorm.DB, value interface{}) (err error) {
	//在draft model上创建publishEvent, 而如果是production, 则跳过
	if publish.IsDraftMode(db) && publish.IsPublishableModel(value) {
		scope := db.NewScope(value)
		var sortingPublishEvent = changedSortingPublishEvent{
			Table: scope.TableName(),
		}
		for _, field := range scope.PrimaryFields() {
			sortingPublishEvent.PrimaryKeys = append(sortingPublishEvent.PrimaryKeys, field.DBName)
		}
		var result []byte
		if result, err = json.Marshal(sortingPublishEvent); err == nil {
			//根据条件查找或者创建一个相同条件的记录
			err = db.New().Where("publish_status = ?", publish.DIRTY).Where(map[string]interface{}{
				"name":     "changed_sorting",
				"argument": string(result),
			}).Attrs(map[string]interface{}{ //默认值
				"publish_status": publish.DIRTY,
				"description":    "Changed sort order for " + scope.GetModelStruct().ModelType.Name(),
			}).FirstOrCreate(&publish.PublishEvent{}).Error
		}
	}
	return err
}

func reorderPositions(scope *gorm.Scope) {
	if !scope.HasError() {
		if _, ok := scope.Value.(sortingInterface); ok {
			table := scope.TableName()
			var additionalSQL []string
			var additionalValues []interface{}
			// 如果db中设置了l10n:locale，并且locale不为空（选择了语言), 同时当前model struct包含了l10n(数据库列language_code), 则只对当前语言进行重新排序
			if locale, ok := scope.DB().Get("l10n:locale"); ok && locale.(string) != "" && l10n.IsLocalizable(scope) {
				additionalSQL = append(additionalSQL, "language_code = ?")
				additionalValues = append(additionalValues, locale)
			}
			additionalValues = append(additionalValues, additionalValues...)

			//如果有deleted_at列，则不在重新排序的范围内
			if scope.HasColumn("DeletedAt") {
				additionalSQL = append(additionalSQL, "deleted_at IS NULL")
			}

			var sql string

			if len(additionalSQL) > 0 {
				//从数据库包中，查找所有比当前position小的行的总数+1,就是当前行新的position
				sql = fmt.Sprintf("UPDATE %v SET position = (SELECT COUNT(pos) + 1 FROM (SELECT DISTINCT(position) AS pos FROM %v WHERE %v) AS t2 WHERE t2.pos < %v.position) WHERE %v", table, table, strings.Join(additionalSQL, " AND "), table, strings.Join(additionalSQL, " AND "))
			} else {
				sql = fmt.Sprintf("UPDATE %v SET position = (SELECT COUNT(pos) + 1 FROM (SELECT DISTINCT(position) AS pos FROM %v) AS t2 WHERE t2.pos < %v.position)", table, table, table)
			}

			if scope.NewDB().Exec(sql, additionalValues...).Error == nil {
				// Create Publish Event
				// 如果更新成功，而且当前db的模式为draftModel, 还需要处理_draft表与production表的同步, 对于production， 则不需要处理
				createPublishEvent(scope.DB(), scope.Value)
			}
		}
	}
}

//对于实现了sortingInterface接口的model, 则默认选择position列进行排序
func beforeQuery(scope *gorm.Scope) {
	v := modelValue(scope.Value)
	//实现了sortingDescInterface SortingDesc方法
	if _, ok := v.(sortingDescInterface); ok {
		scope.Search.Order("position desc")
	} else if _, ok := v.(sortingInterface); ok {
		scope.Search.Order("position")
	}
}

func RegisterCallbacks(db *gorm.DB) {
	db.Callback().Create().Before("gorm:create").Register("sorting:initalize_position", initalizePosition)
	//删除一个元素后，重新对其它的行进行排序
	db.Callback().Delete().After("gorm:after_delete").Register("sorting:reorder_positions", reorderPositions)
	db.Callback().Query().Before("gorm:query").Register("sorting:sort_by_position", beforeQuery)
}
