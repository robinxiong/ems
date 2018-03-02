package utils

import (
	"net/url"
	"path"
	"strings"
	"regexp"
)

//handler path, request 路径
//捕获路由中的参数
// source = "/test/:id" pth="/test/123456
//返回匹配到的参数，匹配的路径，以及bool是否配匹配
func ParamsMatch(source string, pth string) (url.Values, string, bool){
	var (
		i, j int
		p    = make(url.Values)
		ext  = path.Ext(pth)
	)
	pth = strings.TrimSuffix(pth, ext)

	if ext != "" {
		p.Add(":format", strings.TrimPrefix(ext, "."))
	}


	for i < len(pth) {
		switch {
		case j >= len(source):
			if source != "/" && len(source) > 0 && source[len(source)-1] == '/' {
				return p, pth[:i], true
			}

			if source == "" && pth == "/" {
				return p, pth, true
			}
			return p, pth[:i], false
		case source[j] == ':':
			var name, val string
			var nextc byte
			//找到数字或者_,-,!
			name, nextc, j = match(source, isAlumn, j+1)
			//找到参数所对应的值，直到'/'或者name上面中的字符
			val, _, i = match(pth, matchPart(nextc), i)

			//如果带有正则表达式 ：id[\d],则在val中匹配到正则表达式中的值
			if (j < len(source)) && source[j] == '[' {
				var index int
				if idx := strings.Index(source[j:], "]/"); idx > 0 {
					index = idx
				} else if source[len(source)-1] == ']' {
					index = len(source) - j - 1
				}

				if index > 0 {
					match := strings.TrimSuffix(strings.TrimPrefix(source[j:j+index+1], "["), "]")
					if reg, err := regexp.Compile("^" + match + "$"); err == nil && reg.MatchString(val) {
						j = j + index + 1
					} else {
						return nil, "", false
					}
				}
			}

			p.Add(":"+name, val)
		case pth[i] == source[j]:
			i++
			j++
		default:
			return nil, "", false
		}
	}

	if j != len(source) {
		if (len(source) == j+1) && source[j] == '/' {
			return p, pth, true
		}

		return nil, "", false
	}
	return p, pth, true
}

func isAlumn(ch byte) bool {
	return isAlpha(ch) || isDigit(ch)
}

func isAlpha(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch == '-' || ch == '!'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

//匹配s中，符合f函数的字节，如果没有找到，则返回子字符串，以及没有匹配的到字符和位置
func match(s string, f func(byte)bool, i int) (matched string, next byte, j int) {
	j = i
	for j < len(s) && f(s[j]) {
		j++
	}
	if j < len(s) {
		next = s[j]
	}
	return s[i:j], next, j
}

func matchPart(b byte) func(byte) bool {
	return func (c byte) bool {
		return c != b && c != '/'
	}
}