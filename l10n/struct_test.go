package l10n

import (
	"ems/test/utils"

	"github.com/jinzhu/gorm"
)

type Product struct {
	ID int `gorm:"primary_key"`
	Locale
}
type ColorVariation struct {
	ID       int `gorm:"primary_key"`
	Quantity int
	Color    Color
}

type Color struct {
	ID   int `gorm:"primary_key"`
	Code string
	Name string
	Locale
}

type Brand struct {
	ID   int `gorm:"primary_key"`
	Name string
	Locale
}

type Tag struct {
	ID   int `gorm:"primary_key"`
	Name string
	Locale
}

type Category struct {
	ID   int `gorm:"primary_key"`
	Name string
	Locale
}

var dbGlobal, dbCN, dbEN *gorm.DB

func init() {
	db := utils.TestDB()
	RegisterCallbacks(db)
	db.DropTableIfExists(&Product{})
	db.DropTableIfExists(&Brand{})
	db.DropTableIfExists(&Tag{})
	db.DropTableIfExists(&Category{})
	db.Exec("drop table product_tags;")
	db.Exec("drop table product_categories;")
	db.AutoMigrate(&Product{}, &Brand{}, &Tag{}, &Category{})  //根据struct创建数据库表， 每一张表都会带有language_code, 它定义在l10n.go
	dbGlobal = db
	dbCN = dbGlobal.Set("l10n:locale", "zh")
	dbEN = dbGlobal.Set("l10n:locale", "en")
}
