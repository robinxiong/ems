package auth_identity

import "time"

type SignLogs struct {
	Log         string `sql:"-"`
	SignInCount uint
	Logs        []SignLog
}
// SignLog sign log
type SignLog struct {
	UserAgent string
	At        *time.Time
	IP        string
}
