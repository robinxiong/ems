package models

import (
	"github.com/jinzhu/gorm"
	"time"
	"github.com/qor/transition"
)

type Order struct {
	gorm.Model
	UserID            uint
	User              User
	PaymentAmount     float32
	AbandonedReason   string
	DiscountValue     uint
	TrackingNumber    *string
	ShippedAt         *time.Time
	ReturnedAt        *time.Time
	CancelledAt       *time.Time
	ShippingAddressID uint `form:"shippingaddress"`
	ShippingAddress   Address
	BillingAddressID  uint `form:"billingaddress"`
	BillingAddress    Address
	OrderItems        []OrderItem
	transition.Transition
}

type OrderItem struct {
	gorm.Model
	OrderID         uint
	SizeVariationID uint `cartitem:"SizeVariationID"`
	SizeVariation   SizeVariation
	Quantity        uint `cartitem:"Quantity"`
	Price           float32
	DiscountRate    uint
	transition.Transition
}
