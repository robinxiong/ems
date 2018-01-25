package sorting

import (
	"github.com/jinzhu/gorm"
	"ems/publish"
)

type Product struct {
	gorm.Model
	Name string
	Sorting
	publish.Status
}


