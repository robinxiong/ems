package admin

import "ems/core"

type Auth interface {
	GetCurrentUser(*Context) core.CurrentUser
	LoginURL(*Context) string
	LogoutURL(*Context) string
}
