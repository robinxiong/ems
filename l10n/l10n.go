package l10n

// Global global language, 没有在gorm.DB中设置l10n:locale时返回
var Global = "en-US"
type l10nInterface interface{
	IsGlobal() bool
	SetLocale(locale string)
}

type localeCreatableInterface interface {
	CreatableFromLocale()
}
type localeCreatableInterface2 interface {
	LocaleCreatable()
}

// Locale embed this struct into GROM-backend models to enable localization feature for your model
type Locale struct {
	LanguageCode string `sql:"size:20" gorm:"primary_key"`
}


// IsGlobal 返回当前locale是否为全局设置
func (l Locale) IsGlobal() bool {
	return l.LanguageCode == Global
}

// 设置Locale或其它包含Locale的struct的LanguageCode值
func (l Locale) SetLocale(locale string) {
	l.LanguageCode = locale
}


// LocalCreatable允许你嵌入到你的model中，它使用数据库资源可以从locales中创建，默认情况下，它仅能从global中创建en-US
// LocaleCreatable if you embed it into your model, it will make the resource be creatable from locales, by default, you can only create it from global
type LocaleCreatable struct {
	Locale
}

func (LocaleCreatable) CreatableFromLocale(){}