package db

import "github.com/jinzhu/gorm"

var (
	DB *gorm.DB
)

func init(){
	var err error
}