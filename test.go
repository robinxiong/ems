package main

import (
	"github.com/jinzhu/gorm"
	"fmt"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
)

type User struct {
	gorm.Model
	Name string
	Password string
}

func main(){
	user := &User{}
	DB, err := gorm.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=True&loc=Local", "root", "root", "127.0.0.1", 3306, "ems"))
	if err != nil {
		log.Print(err)
	}
	isExist := DB.HasTable(&user)
	log.Print(isExist)
	typeOf := reflect.TypeOf(user)
	valueOf := reflect.ValueOf(user)
	log.Print(typeOf.Kind(), valueOf.Kind())
}