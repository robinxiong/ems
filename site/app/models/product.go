package models

import (
	"github.com/jinzhu/gorm"
	"ems/slug"
	"ems/seo"
	"ems/publish2"
	"ems/media/media_library"
	"ems/l10n"
	"ems/sorting"
)

type Product struct {
	gorm.Model
	l10n.Locale
	sorting.SortingDESC

	Name                  string
	NameWithSlug          slug.Slug    `l10n:"sync"`
	Code                  string       `l10n:"sync"`
	CategoryID            uint         `l10n:"sync"`
	Category              Category     `l10n:"sync"`
	Collections           []Collection `l10n:"sync" gorm:"many2many:product_collections;"`
	MadeCountry           string       `l10n:"sync"`
	Gender                string       `l10n:"sync"`
	MainImage             media_library.MediaBox
	Price                 float32          `l10n:"sync"`
	Description           string           `sql:"size:2000"`
	ColorVariations       []ColorVariation `l10n:"sync"`
	ColorVariationsSorter sorting.SortableCollection
	ProductProperties     ProductProperties `sql:"type:text"`
	Seo                   seo.Setting

	Variations []ProductVariation

	publish2.Version
	publish2.Schedule
	publish2.Visible
}

type ProductVariation struct {
	gorm.Model
	ProductID *uint
	Product   Product

	Color      Color `variations:"primary"`
	ColorID    *uint
	Size       Size `variations:"primary"`
	SizeID     *uint
	Material   Material `variations:"primary"`
	MaterialID *uint

	SKU               string
	ReceiptName       string
	Featured          bool
	Price             uint
	SellingPrice      uint
	AvailableQuantity uint
	Images            media_library.MediaBox
}


type ColorVariation struct {
	gorm.Model
	ProductID      uint
	Product        Product
	ColorID        uint
	Color          Color
	ColorCode      string
	Images         media_library.MediaBox
	SizeVariations []SizeVariation
	publish2.SharedVersion
}


type ProductProperties []ProductProperty

type ProductProperty struct {
	Name  string
	Value string
}

type SizeVariation struct {
	gorm.Model
	ColorVariationID  uint
	ColorVariation    ColorVariation
	SizeID            uint
	Size              Size
	AvailableQuantity uint
	publish2.SharedVersion
}

type ProductImage struct {
	gorm.Model
	Title        string
	Color        Color
	ColorID      uint
	Category     Category
	CategoryID   uint
	SelectedType string
	File         media_library.MediaLibraryStorage `sql:"size:4294967295;" media_library:"url:/system/{{class}}/{{primary_key}}/{{column}}.{{extension}}"`
}