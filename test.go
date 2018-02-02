package main

import (
	"reflect"

	"github.com/jinzhu/gorm"
	"log"
)

type Product struct {
	gorm.Model
	Name string
}

func main() {
	var a float64 = 3.4
	v := reflect.ValueOf(&a).Elem()
	log.Println(v.CanAddr())

	p := &Product{Name:"Robin"}
	v1 := reflect.ValueOf(&p)
	log.Println(v1.CanSet())
}
