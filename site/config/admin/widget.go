package admin

import (
	"ems/widget"
	"ems/l10n"
)

type QorWidgetSetting struct {
	widget.QorWidgetSetting
	// publish2.Version
	// publish2.Schedule
	// publish2.Visible
	l10n.Locale
}