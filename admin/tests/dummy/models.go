package dummy

import (
	"github.com/jinzhu/gorm"
	"time"
	"ems/media/oss"
)

type User struct {
	gorm.Model
	Name         string `gorm:"size:50"`
	Age          uint
	Role         string
	Active       bool
	RegisteredAt *time.Time
	Profile      Profile // has one
	CreditCardID uint
	CreditCard   CreditCard // belongs to
	Addresses    []Address  // has many
	CompanyID    uint
	Company      *Company   // belongs to
	Languages    []Language `gorm:"many2many:user_languages;"` // many 2 many
	Avatar oss.OSS
}


type CreditCard struct {
	gorm.Model
	Number string
	Issuer string
}

type Company struct {
	gorm.Model
	Name string
}

type Address struct {
	gorm.Model
	UserID   uint
	Address1 string
	Address2 string
}

type Language struct {
	gorm.Model
	Name string
}


type Profile struct {
	gorm.Model
	UserID uint
	Name   string
	Sex    string

	Phone Phone
}

type Phone struct {
	gorm.Model

	ProfileID uint64
	Num       string
}
