package l10n
// Locale embed this struct into GROM-backend models to enable localization feature for your model
type Locale struct {
	LanguageCode string `sql:"size:20" gorm:"primary_key"`
}

// LocaleCreatable if you embed it into your model, it will make the resource be creatable from locales, by default, you can only create it from global
type LocaleCreatable struct {
	Locale
}