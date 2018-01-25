package l10n

import (
	"reflect"

	"ems/core/utils"

	"github.com/jinzhu/gorm"
)

func getInterface(scope *gorm.Scope) (i interface{}) {
	if scope.Value == nil {
		return
	}
	value := reflect.ValueOf(scope.Value)
	reflectType := value.Type()
	if reflectType.Kind() == reflect.Slice || reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
	}

	if reflectType.Kind() != reflect.Struct {
		return
	}
	return value.Interface()
}

//当一个数据库表model, 在创建时或者更新时，有一个数据库操作会话会传入beforeCreate/beforeUpdate, 通过这个会话可以获得model要保存或者更新数所库时的值
func IsLocalizable(scope *gorm.Scope) (IsLocalizable bool) {

	if scope.Value == nil {
		return
	}

	_, IsLocalizable = getInterface(scope).(l10nInterface)

	return
	/*
		原始方法
		if scope.GetModelStruct().ModelType == nil {
			return false
		}
		_, IsLocalizable = reflect.New(scope.GetModelStruct().ModelType).Interface().(l10nInterface)
		return
	*/
}

// 确认struct是否满足localeCreatableInterface接口或者localeCreatableInterface2, 如果都没有实现BeforeCreate会检查你的model是否有设置主键，通常locale对像都是设置了主键
func isLocaleCreatable(scope *gorm.Scope) (ok bool) {
	i := getInterface(scope)
	if _, ok = i.(localeCreatableInterface); ok {
		return
	}
	_, ok = i.(localeCreatableInterface2)
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
	/*
	i := getInterface(scope)
	if v, ok := i.(l10nInterface); ok {
		//log.Println(v, locale)
		v.SetLocale(locale) //直接修改了底层数据是无效的，因为在创建对像是，会调用scope.Fields()，检查这个field是否为空，除非调用isblank（检查field.Field)是否为空
	}
	*/
	for _, field := range scope.Fields() {
		if field.Name == "LanguageCode" {
			field.Set(locale) //222 ns/op
		}
	}
	//scope.SetColumn("LanguageCode", locale) 1000/ns
	/*
		// copeSuite.BenchmarkSetLocale    5000000               572 ns/op 每次测试循环所花时间
		v, _ := reflect.New(scope.GetModelStruct().ModelType).Interface().(l10nInterface)
		v.SetLocale(locale)
	*/
}

func isSyncField(field *gorm.StructField) bool {
	if _, ok := utils.ParseTagOption(field.Tag.Get("l10n"))["SYNC"]; ok {
		return true
	}
	return false
}

//如果struct field设置了l10n:sync字段，表示这些字段在更新时保持不变,即添加到db.Search.Omits中
//可以参考callbacks.beforeUpdate
func syncColumns(scope *gorm.Scope) (columns []string) {
	for _, field := range scope.GetModelStruct().StructFields {
		if isSyncField(field) {
			columns = append(columns, field.DBName) //field.DBName 字段在数据库中的名称
		}
	}
	return
}
