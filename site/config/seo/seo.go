package seo

import (
	"ems/seo"
	"ems/l10n"
)

type MySEOSetting struct {
	seo.SystemSEOSetting
	l10n.Locale
}

type SEOGlobalSetting struct {
	SiteName string
}

var SEOCollection *seo.Collection
