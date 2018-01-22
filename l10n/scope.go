package l10n

import (
	"github.com/jinzhu/gorm"
	"reflect"
)

//当一个数据库表model, 在创建时或者更新时，有一个数据库操作会话会传入beforeCreate/beforeUpdate, 通过这个会话可以获得model要保存或者更新数所库时的值

func IsLocalizable(scope *gorm.Scope) (IsLocalizable bool){
	if scope.GetModelStruct().ModelType == nil {
		return false
	}
	_, IsLocalizable = reflect.New(scope.GetModelStruct().ModelType).Interface().(l10nInterface)
	return
}


func getQueryLocale(scope *gorm.Scope) (locale string, isLocale bool) {
	if str, ok := scope.DB().Get("l10n:locale"); ok {
		if locale, ok := str.(string); ok && locale != "" {
			return locale, locale != Global
		}
	}
	return Global, false
}

//每一个gorm.DB中设置的全局变量，e.g.: db.Set("l10n:locale", "zh-CN")
func getLocale(scope *gorm.Scope)(locale string, isLocale bool){
	if str, ok := scope.DB().Get("l10n:localize_to"); ok {
		if locale, ok := str.(string); ok && locale != "" {
			return locale, locale != Global
		}
	}
	return getQueryLocale(scope)
}
