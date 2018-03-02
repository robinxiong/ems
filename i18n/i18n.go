package i18n

import (
	"ems/cache"
	"ems/cache/memory"
	"errors"
	"html/template"
	"strings"
	"github.com/theplant/cldr"
)

var Default = "en-US"

// Translation 一箱翻译的数据结果，包含key, locale两个主键，以及Value表示所对应的值
type Translation struct {
	Key     string
	Locale  string
	Value   string
	Backend Backend `json:"-"`
}

// Backend用来保存翻译数据的方式，可以是数据库表，也可以是yaml文件
type Backend interface {
	LoadTranslations() []*Translation
	SaveTranslation(*Translation) error
	DeleteTranslation(*Translation) error
}

// I18n struct 保存所有的翻译
type I18n struct {
	value      string                    //翻译的值
	Backends   []Backend                 //可以是多个保存方式, 数据库或者yaml文件
	cacheStore cache.CacheStoreInterface //使用哪种方式来缓存翻译的结果
}

func New(backends ...Backend) *I18n {
	i18n := &I18n{Backends: backends, cacheStore: memory.New()}
	i18n.loadToCacheStore() //将加载到的翻译数据，保存到memory cache中
	return i18n
}

func (i18n *I18n) loadToCacheStore() {
	backends := i18n.Backends
	for i := len(backends) - 1; i >= 0; i-- {
		var backend = backends[i]
		for _, translation := range backend.LoadTranslations() {
			i18n.AddTranslation(translation)
		}
	}
}

// AddTranslation add translation
func (i18n *I18n) AddTranslation(translation *Translation) error {
	key := cacheKey(translation.Locale, translation.Key)
	return i18n.cacheStore.Set(key, translation)
}

//translation.Locale, translation.Key
func cacheKey(strs ...string) string {
	return strings.Join(strs, "/")
}
func (i18n *I18n) SaveTranslation(translation *Translation) error {
	for _, backend := range i18n.Backends {
		if backend.SaveTranslation(translation) == nil {
			i18n.AddTranslation(translation)
			return nil
		}
	}

	return errors.New("failed to save translation")
}
func (i18n *I18n) T(locale, key string, args ...interface{}) template.HTML {
	var (
		value          = i18n.value
		translationKey = key
	)

	if locale == "" {
		locale = Default //如果没有指定locale(cookie中没有locale的值)，则使用默认的en-US
	}

	var translation Translation

	if err := i18n.cacheStore.Unmarshal(cacheKey(locale, key), &translation); err != nil || translation.Value == "" {
		//如果没有在locale语言中找到值，则从默认的en-US中查找
		if translation.Value == "" {
			if err := i18n.cacheStore.Unmarshal(cacheKey(Default, key), &translation); err != nil || translation.Value == "" {
				//默认的值也没有找到，将它保存到数据库中
				var defaultBackend Backend
				if len(i18n.Backends) > 0 {
					defaultBackend = i18n.Backends[0]
				}
				translation = Translation{Key: translationKey, Value: value, Locale: locale, Backend: defaultBackend}

				// Save translation
				i18n.SaveTranslation(&translation)
			}
		}
	}

	if translation.Value != "" {
		value = translation.Value
	} else {
		//locale以及default locale都没有找到key所对应的值，则直接显示key
		value = key
	}

	//调用theplant公共包，翻译常用的词，比如1月，2月，星期一，星期二等
	if str, err := cldr.Parse(locale, value, args...); err == nil {
		value = str
	}

	return template.HTML(value)
}
