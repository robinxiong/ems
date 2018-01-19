package notification

import "ems/admin"


type ActionArgument struct {
	Message  *NotificationMessage
	Context  *admin.Context
	Argument interface{}
}

type Action struct {
	Name         string
	Label        string
	Method       string
	MessageTypes []string
	Resource     *admin.Resource
	Visible      func(data *NotificationMessage, context *admin.Context) bool
	URL          func(data *NotificationMessage, context *admin.Context) string
	Handler      func(actionArgument *ActionArgument) error
	Undo         func(actionArgument *ActionArgument) error
	FlashMessage func(actionArgument *ActionArgument, succeed bool, isUndo bool) string
}
