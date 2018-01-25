package main

import (
	"ems/test/utils"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Product struct {
	ID         int        `gorm:"primary_key"`
	Categories []Category `gorm:"many2many:product_categories;ForeignKey:id;AssociationForeignKey:id"`
}
type Category struct {
	ID   int `gorm:"primary_key"`
	Name string
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	db := utils.TestDB()
	db.DropTableIfExists(&Product{})
	db.DropTableIfExists(&Category{})
	db.Exec("drop table product_categories;")
	//创建producs, categories表
	db.AutoMigrate(&Product{}, &Category{})
}
