package main

import (

	"ems/publish"
	"ems/test/utils"
	"log"
	"github.com/jinzhu/gorm"
)



func main() {
	db := utils.TestDB()
	pb := publish.New(db)
	log.Println(pb)
	log.Println(gorm.DefaultTableNameHandler(db, "hello"))
}
