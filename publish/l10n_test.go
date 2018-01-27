package publish

import (
	"github.com/jinzhu/gorm"
	"ems/l10n"
)

type Book struct {
	gorm.Model
	l10n.Locale
	Status
	Name        string
	CategoryID  uint
	Category    Category
	PublisherID uint
	Publisher   Publisher
	Comments    []Comment
	Authors     []Author `gorm:"many2many:author_books;ForeignKey:ID;AssociationForeignKey:ID"`
}


type Publisher struct {
	gorm.Model
	Status
	Name string
}

type Comment struct {
	gorm.Model
	l10n.Locale
	Status
	Content string
	BookID  uint
}

type Author struct {
	gorm.Model
	l10n.Locale
	Name string
}