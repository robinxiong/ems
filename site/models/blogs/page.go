package blogs

import (
	"ems/publish2"
)

type Page struct {

	publish2.Version
	publish2.Schedule
	publish2.Visible
}
