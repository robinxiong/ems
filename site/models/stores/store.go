package stores

import (
	"github.com/jinzhu/gorm"
	"ems/sorting"
	"ems/location"
)

type Store struct {
	gorm.Model
	StoreName string
	Owner     Owner
	Phone     string
	Email     string
	location.Location
	sorting.Sorting
}


type Owner struct {
	Name    string
	Contact string
	Email   string
}