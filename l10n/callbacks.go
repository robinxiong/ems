package l10n

import (
	"github.com/jinzhu/gorm"
)

func beforeCreate(scope *gorm.Scope) {

	if IsLocalizable(scope) {

		if locale, ok := getLocale(scope); ok {
			if isLocaleCreatable(scope) || !scope.PrimaryKeyZero() {
				setLocale(scope, locale)
			}
		}
	}
}

func RegisterCallbacks(db *gorm.DB) {
	callback := db.Callback()
	callback.Create().Before("gorm:before_create").Register("l10n:before_create", beforeCreate)
}
