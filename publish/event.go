package publish

import "github.com/jinzhu/gorm"

// EventInterface defined methods needs for a publish event
type EventInterface interface {
	Publish(db *gorm.DB, event PublishEventInterface) error
	Discard(db *gorm.DB, event PublishEventInterface) error
}

var events = map[string]EventInterface{}

// RegisterEvent register publish event
func RegisterEvent(name string, event EventInterface) {
	events[name] = event
}

// PublishEvent default publish event model
type PublishEvent struct {
	gorm.Model
	Name          string
	Description   string
	Argument      string `sql:"size:65532"`
	PublishStatus bool
	PublishedBy   string
}