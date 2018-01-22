package l10n

// Global global language
var Global = "en-US"
type l10nInterface interface{
	IsGlobal() bool
	SetLocal(locale string)
}
// Locale embed this struct into GROM-backend models to enable localization feature for your model
type Locale struct {
	LanguageCode string `sql:"size:20" gorm:"primary_key"`
}


//IsGlobal 返回当前locale是否为英文
func (l Locale) IsGlobal() bool {
	return l.LanguageCode == Global
}

func (l Locale) SetLocale(locale string) {
	l.LanguageCode = locale
}


// LocaleCreatable if you embed it into your model, it will make the resource be creatable from locales, by default, you can only create it from global
type LocaleCreatable struct {
	Locale
}
