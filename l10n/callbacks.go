package l10n

import (
	"fmt"

	"ems/core/utils"
	"reflect"

	"github.com/jinzhu/gorm"
)

func beforeQuery(scope *gorm.Scope) {

	if IsLocalizable(scope) {
		quotedTableName := scope.QuotedTableName()
		quotedPrimaryKey := scope.Quote(scope.PrimaryKey()) //id column or first primary key
		_, hasDeletedAtColumn := scope.FieldByName("deleted_at")
		locale, isLocale := getQueryLocale(scope) //db中是否设置了l10n:locale

		//先确认db中是否设置了l10n的模式，即db.Set("10n:mode", "unscoped")
		switch mode, _ := scope.DB().Get("l10n:mode"); mode {
		case "unscoped":
			//不指定作何language_code语句,  比如在afterUpdate之后，需要更新sync字段
		case "global":
			scope.Search.Where(fmt.Sprintf("%v.language_code = ?", quotedTableName), Global)
		case "locale":
			//sorting/callbacks.go initalizePosition
			//如果l10n:locale为空，没有设置，则返回所有locale, zh-CN, en-US

			scope.Search.Where(fmt.Sprintf("%v.language_code = ?", quotedTableName), locale)
			//取反，查询不在当前locale的全局行
		case "reverse":
			if !scope.Search.Unscoped && hasDeletedAtColumn {
				//如果数所库表有deleted_at列, 而查找不在（language=local同时又没有删除的行），同时又有language=global的行
				//即查询当前没有包含zh-CN的global行
				scope.Search.Where(fmt.Sprintf(
					"(%v.%v NOT IN (SELECT DISTINCT(%v) FROM %v t2 WHERE t2.language_code = ? AND t2.deleted_at IS NULL) AND %v.language_code = ?)", quotedTableName, quotedPrimaryKey, quotedPrimaryKey, quotedTableName, quotedTableName), locale, Global)
			} else {
				scope.Search.Where(fmt.Sprintf("(%v.%v NOT IN (SELECT DISTINCT(%v) FROM %v t2 WHERE t2.language_code = ?) AND %v.language_code = ?)", quotedTableName, quotedPrimaryKey, quotedPrimaryKey, quotedTableName, quotedTableName), locale, Global)
			}
		case "fallback":
			fallthrough
		default:
			if isLocale {
				if !scope.Search.Unscoped && hasDeletedAtColumn {
					scope.Search.Where(fmt.Sprintf("((%v.%v NOT IN (SELECT DISTINCT(%v) FROM %v t2 WHERE t2.language_code = ? AND t2.deleted_at IS NULL) AND %v.language_code = ?) OR %v.language_code = ?) AND %v.deleted_at IS NULL", quotedTableName, quotedPrimaryKey, quotedPrimaryKey, quotedTableName, quotedTableName, quotedTableName, quotedTableName), locale, Global, locale)
				} else {
					scope.Search.Where(fmt.Sprintf("(%v.%v NOT IN (SELECT DISTINCT(%v) FROM %v t2 WHERE t2.language_code = ?) AND %v.language_code = ?) OR (%v.language_code = ?)", quotedTableName, quotedPrimaryKey, quotedPrimaryKey, quotedTableName, quotedTableName, quotedTableName), locale, Global, locale)
				}
				scope.Search.Order(gorm.Expr(fmt.Sprintf("%v.language_code = ? DESC", quotedTableName), locale))
			} else {
				scope.Search.Where(fmt.Sprintf("%v.language_code = ?", quotedTableName), Global)
			}
		}

	}
}

func beforeUpdate(scope *gorm.Scope) {
	if IsLocalizable(scope) {
		locale, isLocale := getLocale(scope)

		switch mode, _ := scope.DB().Get("l10n:mode"); mode {
		case "unscoped":
		default:
			//设置为global的值, 如果更新成功，在调用afterUpdate更新其它l10n:sync的columns
			scope.Search.Where(fmt.Sprintf("%v.language_code = ?", scope.QuotedTableName()), locale)
			setLocale(scope, locale) //设置scope.value的languare的值
		}
		//如果是locale，则忽略l10n:sync的列，如果需要更新，则通过dbGlobal来更新
		//同时使用UpdateColumns方法
		if isLocale {
			columns := syncColumns(scope)
			scope.Search.Omit(columns...)
		}
	}
}

func afterUpdate(scope *gorm.Scope) {
	if !scope.HasError() {
		//包含了Locale struct
		if IsLocalizable(scope) {
			//如果在db中设置了l10n:locale
			if locale, ok := getLocale(scope); ok {
				//没有任何更新同时要更新的对像主键为非零值,
				if scope.DB().RowsAffected == 0 && !scope.PrimaryKeyZero() {
					var count int
					var query = fmt.Sprintf("%v.language_code = ? AND %v.%v = ?", scope.QuotedTableName(), scope.QuotedTableName(), scope.PrimaryKey())
					//如果包含了delete_at列, 而且它的值为非null，则删除当前本地化行，比如zh-CN
					//将删除后的product重新恢复
					if scope.HasColumn("DeletedAt") {
						scope.NewDB().Unscoped().Where("deleted_at is not null").Where(query, locale, scope.PrimaryKeyValue()).Delete(scope.Value)
					}
					// if no localized records exist, localize it
					if scope.NewDB().Table(scope.TableName()).Where(query, locale, scope.PrimaryKeyValue()).Count(&count); count == 0 {
						scope.DB().RowsAffected = scope.DB().Create(scope.Value).RowsAffected
					}
				}
			} else if syncColumns := syncColumns(scope); len(syncColumns) > 0 { // is global, 需要更新l10n:sync字段
				if mode, _ := scope.DB().Get("l10n:mode"); mode != "unscoped" {
					if scope.DB().RowsAffected > 0 {
						var primaryField = scope.PrimaryField()
						var syncAttrs = map[string]interface{}{}
						//gorm/callback_update.go scope.InstanceSet("gorm:update_attrs", updateMaps)
						//可以参考gorm/main.go func (s *DB) Updates(values interface{}, ignoreProtectedAttrs ...bool) *DB
						//获取Updates或者UpdateColumns方法中指定的map[string]string{}要更新的列, 同时它没有在scope.Search.Omit中被忽略
						if updateAttrs, ok := scope.InstanceGet("gorm:update_attrs"); ok {
							for key, value := range updateAttrs.(map[string]interface{}) {
								for _, syncColumn := range syncColumns {
									if syncColumn == key {
										syncAttrs[syncColumn] = value
										break
									}
								}
							}
						} else {
							for _, syncColumn := range syncColumns {
								if field, ok := scope.FieldByName(syncColumn); ok && field.IsNormal {
									syncAttrs[syncColumn] = field.Field.Interface()
								}
							}
						}

						if len(syncAttrs) > 0 {
							db := scope.DB().Model(reflect.New(utils.ModelType(scope.Value)).Interface()).Set("l10n:mode", "unscoped").Where("language_code <> ?", Global)
							if !primaryField.IsBlank {
								db = db.Where(fmt.Sprintf("%v = ?", primaryField.DBName), primaryField.Field.Interface())
							}
							scope.Err(db.UpdateColumns(syncAttrs).Error)
						}
					}
				}
			}
		}
	}
}
func beforeCreate(scope *gorm.Scope) {

	if IsLocalizable(scope) {
		//是否有在db中设置l10n:locale, zh-CN, 如果没有，则设置为全局的值en-US  from language_code
		//如果isLocaleCreatable, 又没有en-US记录，则不能创建zh-CN, 及en, 可以保证需要先创建一条全局的记录
		// 1 en-us 必须先创建全局, 创建完成后，返回 ID:1,language_code:en-us, 如果此记录创建失败，则报错，并且不在继续处理
		// 1 zh
		// 1 en

		//非全局设置

		if locale, ok := getLocale(scope); ok {
			//PrimaryKeyZero用来判断第一个主键的值是否为零值, 这里为false
			if isLocaleCreatable(scope) || !scope.PrimaryKeyZero() {

				setLocale(scope, locale)
				//i, _ := scope.Value.(l10nInterface)

			} else {
				err := fmt.Errorf("the resource %v cannot be created in %v", scope.GetModelStruct().ModelType.Name(), locale)
				scope.Err(err) //报错，并且不在继续处理
			}
		} else {

			setLocale(scope, Global)

		}

	}
}
func beforeDelete(scope *gorm.Scope) {
	if IsLocalizable(scope) {
		if locale, ok := getQueryLocale(scope); ok { // is locale
			scope.Search.Where(fmt.Sprintf("%v.language_code = ?", scope.QuotedTableName()), locale)
		}
	}
}
func RegisterCallbacks(db *gorm.DB) {
	callback := db.Callback()
	callback.Create().Before("gorm:before_create").Register("l10n:before_create", beforeCreate)
	callback.RowQuery().Before("gorm:row_query").Register("l10n:before_query", beforeQuery)
	callback.Query().Before("gorm:query").Register("l10n:before_query", beforeQuery)
	callback.Update().Before("gorm:before_update").Register("l10n:before_update", beforeUpdate)
	callback.Update().After("gorm:after_update").Register("l10n:after_update", afterUpdate)
	callback.Delete().Before("gorm:before_delete").Register("l10n:before_delete", beforeDelete)
}
