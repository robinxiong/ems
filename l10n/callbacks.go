package l10n

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

func beforeCreate(scope *gorm.Scope) {

	if IsLocalizable(scope) {
		//是否有在db中设置l10n:locale, zh-CN, 如果没有，则设置为全局的值en-US  from language_code
		//如果isLocaleCreatable, 又没有en-US记录，则不能创建zh-CN, 及en, 可以保证需要先创建一条全局的记录
		// 1 en-us 必须先创建全局, 创建完成后，返回 ID:1,language_code:en-us, 如果此记录创建失败，则报错，并且不在继续处理
		// 1 zh
		// 1 en
		if locale, ok := getLocale(scope); ok {
			//PrimaryKeyZero用来判断第一个主键的值是否为零值, 这里为false
			if isLocaleCreatable(scope) || !scope.PrimaryKeyZero() {
				setLocale(scope, locale)
			} else {
				err := fmt.Errorf("the resource %v cannot be created in %v", scope.GetModelStruct().ModelType.Name(), locale)
				scope.Err(err) //报错，并且不在继续处理
			}
		} else {
			setLocale(scope, Global)
		}
	}
}

func RegisterCallbacks(db *gorm.DB) {
	callback := db.Callback()
	callback.Create().Before("gorm:before_create").Register("l10n:before_create", beforeCreate)
}
