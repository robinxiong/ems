package models

import (
	"github.com/jinzhu/gorm"
	"time"
	"ems/media/oss"
)

type User struct {
	gorm.Model
	Email                  string `form:"email"`
	Password               string
	Name                   string `form:"name"`
	Gender                 string
	Role                   string
	Birthday               *time.Time
	Balance                float32
	DefaultBillingAddress  uint `form:"default-billing-address"`
	DefaultShippingAddress uint `form:"default-shipping-address"`
	Addresses              []Address
	Orders                 []Order
	Avatar                 AvatarImageStorage

	// Confirm
	ConfirmToken string
	Confirmed    bool

	// Recover
	RecoverToken       string
	RecoverTokenExpiry *time.Time

	// Accepts
	AcceptPrivate bool `form:"accept-private"`
	AcceptLicense bool `form:"accept-license"`
	AcceptNews    bool `form:"accept-news"`
}

//实现core.CurrentUser接口，否则config/auth/admin_auth.go将返回错误
func (user User) DisplayName() string {
	return user.Email
}



type AvatarImageStorage struct{ oss.OSS }