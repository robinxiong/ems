package sorting

import (
	"github.com/jinzhu/gorm"
	"ems/l10n"
)

type Brand struct {
	gorm.Model
	l10n.Locale
	Sorting
	Name string
}