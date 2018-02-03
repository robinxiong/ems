package models

import (
	"github.com/qor/page_builder"
	"ems/publish2"
)

type Page struct {
	page_builder.Page

	publish2.Version
	publish2.Schedule
	publish2.Visible
}