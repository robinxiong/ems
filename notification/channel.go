package notification

import (
	"ems/core"
)

func (notification *Notification) RegisterChannel(channel ChannelInterface) {
	notification.Channels = append(notification.Channels, channel)
}

type ChannelInterface interface {
	Send(message *Message, context *core.Context) error
	GetNotifications(user interface{}, results *NotificationsResult, notification *Notification, context *core.Context) error
	GetUnresolvedNotificationsCount(user interface{}, notification *Notification, context *core.Context) uint
	GetNotification(user interface{}, notificationID string, notification *Notification, context *core.Context) (*NotificationMessage, error)
}
