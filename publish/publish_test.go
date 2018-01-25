package publish

import (
	"github.com/jinzhu/gorm"
	"ems/test/utils"
	"ems/l10n"
	"testing"
	"github.com/stretchr/testify/assert"
)

var pb *Publish
var pbdraft *gorm.DB
var pbprod *gorm.DB
var db *gorm.DB


type Product struct {
	gorm.Model
	Name       string
	Quantity   uint
	Color      Color
	ColorId    int
	Categories []Category `gorm:"many2many:product_categories"`
	Languages  []Language `gorm:"many2many:product_languages"`
	Status
}

type Color struct {
	gorm.Model
	Name string
}

type Language struct {
	gorm.Model
	Name string
}

type Category struct {
	gorm.Model
	Name string
	Status
}


func init(){
	db = utils.TestDB()
	l10n.RegisterCallbacks(db)
	pb = New(db)
}


func TestIsPublishableModel(t *testing.T) {
	product := Product{}
	color := Color{}
	assert.True(t, IsPublishableModel(&product), "Product is implement publishInterface")
	assert.False(t, IsPublishableModel(&color), "color isn't implement publishInterface, because it not embbed Status struct")
}