package inline_edit

import (
	"ems/i18n"
	"html/template"

)

func FuncMap(I18n *i18n.I18n, locale string, enableInlineEdit bool) template.FuncMap {
	return template.FuncMap{
		"t": InlineEdit(I18n, locale, enableInlineEdit),
	}
}

func InlineEdit(I18n *i18n.I18n, locale string, isInline bool) func(string, ...interface{}) template.HTML {
	return func(key string, args ...interface{}) template.HTML {
		var value template.HTML
		value = I18n.T(locale, key)
		return value
	}
}
