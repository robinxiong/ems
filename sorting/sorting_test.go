package sorting

import (
	"github.com/jinzhu/gorm"
	"ems/test/utils"
	"ems/l10n"
)

type User struct {
	gorm.Model
	Name string
	Sorting
}

var db *gorm.DB

func init(){
	db = utils.TestDB()
	RegisterCallbacks(db)
	l10n.RegisterCallbacks(db)

}