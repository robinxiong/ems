package notification

type Notification struct {
	Config   *Config
	Channels []ChannelInterface
	Actions  []*Action
}

func New(config *Config) *Notification {
	notification := &Notification{Config: config}
	return notification
}


type NotificationsResult struct {
	Notification  *Notification
	Notifications []*NotificationMessage
	Resolved      []*NotificationMessage
}