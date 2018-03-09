package orders

import (
	"github.com/jinzhu/gorm"
	"time"

	"ems/site/models/users"
)

type Order struct {
	gorm.Model
	UserID            uint
	User              users.User
	PaymentAmount     float32
	AbandonedReason   string
	DiscountValue     uint
	TrackingNumber    *string
	ShippedAt         *time.Time
	ReturnedAt        *time.Time
	CancelledAt       *time.Time
	ShippingAddressID uint `form:"shippingaddress"`
	ShippingAddress   users.Address
	BillingAddressID  uint `form:"billingaddress"`
	BillingAddress    users.Address
	OrderItems        []OrderItem

}

type OrderItem struct {
	gorm.Model
	OrderID         uint
	SizeVariationID uint `cartitem:"SizeVariationID"`
	Quantity        uint `cartitem:"Quantity"`
	Price           float32
	DiscountRate    uint

}
