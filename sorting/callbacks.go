package sorting

import (
	"github.com/jinzhu/gorm"
	"reflect"
)

func initalizePosition(scope *gorm.Scope){
	if !scope.HasError() {
		if _, ok := scope.Value.(sortingInterface); ok {
			var lastPosition int
			//
			scope.NewDB().Set("l10n:mode", "locale").Model(modelValue(scope.Value)).Select("position").Order("position DESC").Limit(1).Row().Scan(&lastPosition)
			scope.SetColumn("Position", lastPosition+1)
		}
	}
}

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

func RegisterCallbacks(db *gorm.DB) {
	db.Callback().Create().Before("gorm:create").Register("sorting:initalize_position", initalizePosition)
}
