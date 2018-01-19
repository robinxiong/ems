package models

import (
	"github.com/jinzhu/gorm"
	"ems/l10n"
)

type Collection struct {
	gorm.Model
	Name string
	l10n.LocaleCreatable
}
