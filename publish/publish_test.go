package publish

import (
	"ems/l10n"
	"ems/test/utils"
	"fmt"
	"testing"

	"github.com/jinzhu/gorm"
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

func init() {
	db = utils.TestDB()
	l10n.RegisterCallbacks(db)

	pb := New(db)

	pbdraft = pb.DraftDB()
	pbprod = pb.ProductionDB()
	//删除所有数据库
	for _, table := range []string{"product_categories", "product_categories_draft", "product_languages", "product_languages_draft", "author_books", "author_books_draft"} {
		pbprod.Exec(fmt.Sprintf("drop table %v", table))
	}

	for _, value := range []interface{}{&Product{}, &Color{}, &Category{}, &Language{}, &Book{}, &Publisher{}, &Comment{}, &Author{}} {
		//因为Color类没有实现publishInterface, 所以pbdraft正常删除了colors表, 而不是colors_draft
		if IsPublishableModel(value) {
			pbdraft.DropTable(value)
		}
		pbprod.DropTable(value)

		//调用publish中定义的AutoMigrate, 它会创建_draft表
		pbprod.AutoMigrate(value)
		pb.AutoMigrate(value) //migrate to draft db
	}
}

func TestIsPublishableModel(t *testing.T) {
	/*	product := Product{}
		color := Color{}
		assert.True(t, IsPublishableModel(&product), "Product is implement publishInterface")
		assert.False(t, IsPublishableModel(&color), "color isn't implement publishInterface, because it not embbed Status struct")*/
}
