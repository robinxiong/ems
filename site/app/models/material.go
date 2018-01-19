package models

import (
	"github.com/jinzhu/gorm"
	"ems/l10n"
)

type Material struct {
	gorm.Model
	l10n.Locale
	Name string
	Code string `l10n:"sync"`
}
