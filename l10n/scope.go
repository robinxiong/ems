package l10n

import (
	"reflect"

	"github.com/jinzhu/gorm"
)

//当一个数据库表model, 在创建时或者更新时，有一个数据库操作会话会传入beforeCreate/beforeUpdate, 通过这个会话可以获得model要保存或者更新数所库时的值

func IsLocalizable(scope *gorm.Scope) (IsLocalizable bool) {
	if scope.GetModelStruct().ModelType == nil {

		return false
	}
	_, IsLocalizable = reflect.New(scope.GetModelStruct().ModelType).Interface().(l10nInterface)
	return
}

// 确认struct是否满足localeCreatableInterface接口或者localeCreatableInterface2, 如果都没有实现BeforeCreate会检查你的model是否有设置主键，通常locale对像都是设置了主键
func isLocaleCreatable(scope *gorm.Scope) (ok bool) {
	if _, ok = reflect.New(scope.GetModelStruct().ModelType).Interface().(localeCreatableInterface); ok {
		return
	}
	_, ok = reflect.New(scope.GetModelStruct().ModelType).Interface().(localeCreatableInterface2)
	return
}

// 查询DB中是否设置了l10n:locale, 如果没有找到，则返回全局的local(en-US)， false
func getQueryLocale(scope *gorm.Scope) (locale string, isLocale bool) {
	if str, ok := scope.DB().Get("l10n:locale"); ok {
		if locale, ok := str.(string); ok && locale != "" {
			return locale, locale != Global
		}
	}
	return Global, false
}

// 查找gorm.DB中设置的全局变量，e.g.: db.Set("l10n:locale", "zh-CN")
func getLocale(scope *gorm.Scope) (locale string, isLocale bool) {
	if str, ok := scope.DB().Get("l10n:localize_to"); ok {
		if locale, ok := str.(string); ok && locale != "" {
			return locale, locale != Global
		}
	}
	return getQueryLocale(scope)
}

//设置scope中的所对应的model所对应struct的字段， 如果这个字段为LanguageCode, 则设置它的值
func setLocale(scope *gorm.Scope, locale string) {
	v, _ := reflect.New(scope.GetModelStruct().ModelType).Interface().(l10nInterface)
	v.SetLocale(locale)

	/*for _, field := range scope.Fields() {
		field.Set(locale)  //13732 ns/op
	}*/
	//scope.SetColumn("LanguageCode", locale)
	/*
	// copeSuite.BenchmarkSetLocale    5000000               572 ns/op 每次测试循环所花时间
	v, _ := reflect.New(scope.GetModelStruct().ModelType).Interface().(l10nInterface)
	v.SetLocale(locale)
	*/
}
