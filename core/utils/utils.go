package utils

import (
	"ems/core"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"github.com/jinzhu/gorm"
	"net/url"
	"path"
)

var AppRoot, _ = os.Getwd()

type ContextKey string

var ContextDBName ContextKey = "ContextDB"

// GOPATH return GOPATH from env
func GOPATH() []string {
	paths := strings.Split(os.Getenv("GOPATH"), string(os.PathListSeparator))
	if len(paths) == 0 {
		fmt.Println("GOPATH doesn't exist")
	}
	return paths
}

// GetLocale get locale from request, cookie, after get the locale, will write the locale to the cookie if possible
// Overwrite the default logic with
//     utils.GetLocale = func(context *qor.Context) string {
//         // ....
//     }
func GetLocale(context *core.Context) string {
	if locale := context.Request.Header.Get("Locale"); locale != "" {
		return locale
	}

	if locale := context.Request.URL.Query().Get("locale"); locale != "" {
		if context.Writer != nil {
			context.Request.Header.Set("Locale", locale)
			SetCookie(http.Cookie{Name: "locale", Value: locale, Expires: time.Now().AddDate(1, 0, 0)}, context)
		}
		return locale
	}

	if locale, err := context.Request.Cookie("locale"); err == nil {
		return locale.Value
	}

	return ""
}
func ParseTagOption(str string) map[string]string {
	tags := strings.Split(str, ";")
	setting := map[string]string{}
	for _, value := range tags {
		v := strings.Split(value, ":")
		k := strings.TrimSpace(strings.ToUpper(v[0]))
		if len(v) == 2 {
			setting[k] = v[1]
		} else {
			setting[k] = k
		}
	}
	return setting
}

// ModelType get value's model type
func ModelType(value interface{}) reflect.Type {
	reflectType := reflect.Indirect(reflect.ValueOf(value)).Type()

	for reflectType.Kind() == reflect.Ptr || reflectType.Kind() == reflect.Slice {
		reflectType = reflectType.Elem()
	}

	return reflectType
}

// SetCookie set cookie for context
func SetCookie(cookie http.Cookie, context *core.Context) {
	cookie.HttpOnly = true

	// set https cookie
	if context.Request != nil && context.Request.URL.Scheme == "https" {
		cookie.Secure = true
	}

	// set default path
	if cookie.Path == "" {
		cookie.Path = "/"
	}

	http.SetCookie(context.Writer, &cookie)
}

func HumanizeString(str string) string {
	var human []rune
	for i, l := range str {
		if i > 0 && isUppercase(byte(l)) {
			if (!isUppercase(str[i-1]) && str[i-1] != ' ') || (i+1 < len(str) && !isUppercase(str[i+1]) && str[i+1] != ' ' && str[i-1] != ' ') {
				human = append(human, rune(' '))
			}
		}
		human = append(human, l)
	}
	return strings.Title(string(human))
}

var asicsiiRegexp = regexp.MustCompile("^(\\w|\\s|-|!)*$")

func ToParamString(str string) string {
	if asicsiiRegexp.MatchString(str) {
		return gorm.ToDBName(strings.Replace(str, " ", "_", -1))
	}
	return slug.Make(str)
}

func isUppercase(char byte) bool {
	return 'A' <= char && char <= 'Z'
}

func ExitWithMsg(msg interface{}, value ...interface{}) {
	fmt.Printf("\n"+filenameWithLineNum()+"\n"+fmt.Sprint(msg)+"\n", value...)
	debug.PrintStack()
}

func filenameWithLineNum() string {
	var total = 10
	var results []string
	for i := 2; i < 15; i++ {
		if _, file, line, ok := runtime.Caller(i); ok {
			total--
			results = append(results[:0],
				append(
					[]string{fmt.Sprintf("%v:%v", strings.TrimPrefix(file, os.Getenv("GOPATH")+"src/"), line)},
					results[0:]...)...)

			if total == 0 {
				return strings.Join(results, "\n")
			}
		}
	}
	return ""
}


// GetAbsURL get absolute URL from request, refer: https://stackoverflow.com/questions/6899069/why-are-request-url-host-and-scheme-blank-in-the-development-server
func GetAbsURL(req *http.Request) url.URL {
	var result url.URL

	if req.URL.IsAbs() {
		return *req.URL
	}

	if domain := req.Header.Get("Origin"); domain != "" {
		parseResult, _ := url.Parse(domain)
		result = *parseResult
	}

	result.Parse(req.RequestURI)
	return result
}


func FileServer(dir http.Dir) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := path.Join(string(dir), r.URL.Path)
		if f, err := os.Stat(p); err == nil && !f.IsDir() {
			http.ServeFile(w, r, p)
			return
		}

		http.NotFound(w, r)
	})
}