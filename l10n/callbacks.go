package l10n

import (
	"github.com/jinzhu/gorm"
	"log"
)

func beforeCreate(scope *gorm.Scope) {
	if IsLocalizable(scope) {
		if locale, ok := getLocale(scope); ok {
			log.Print(locale)
		}
	}
}


func RegisterCallbacks(db *gorm.DB) {
	log.Println(db)
}