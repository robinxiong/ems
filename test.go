package main

import (
	"ems/test/utils"
	"log"

	"github.com/jinzhu/gorm"
)

type Product struct {
	gorm.Model
	Name string
}

func main() {
	db := utils.TestDB()
	scope := db.NewScope(&Product{})
	for _, field := range scope.Fields() {
		log.Println(scope.Quote(field.DBName))
	}
}