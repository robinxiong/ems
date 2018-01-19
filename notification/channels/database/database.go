package database

import (
	"ems/core"
	"ems/notification"
	"fmt"

	"github.com/jinzhu/gorm"
)

type Config struct {
	DB *gorm.DB
}

func New(config *Config) *Database {
	if config.DB != nil {
		config.DB.AutoMigrate(&notification.NotificationMessage{})
	} else {
		fmt.Println("Need to have gorm DB in the configuration in order to run migrations")
	}
	return &Database{Config: config}
}

type Database struct {
	Config *Config
}

func (database *Database) Send(message *notification.Message, context *core.Context) error {
	return nil
}

func (database *Database) GetNotifications(user interface{}, results *notification.NotificationsResult, _ *notification.Notification, context *core.Context) error {
	return nil
}

func (database *Database) GetUnresolvedNotificationsCount(user interface{}, notification *notification.Notification, context *core.Context) uint {
	return 0
}
func (database *Database) GetNotification(user interface{}, notificationID string, _ *notification.Notification, context *core.Context) (*notification.NotificationMessage, error) {
	var (
		notice notification.NotificationMessage
	)

	notice = notification.NotificationMessage{}
	return &notice, nil
}
