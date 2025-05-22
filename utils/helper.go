package utils

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type MapSupportedTypes interface {
	string | int64 | float64 | bool
}

// GetMapSpecificValue 获取map的特定类型的值, 相较于GetMapValue, 不用每次反射获取值
func GetMapSpecificValue[T MapSupportedTypes](m map[string]any, key string) T {
	var zero T

	value, exists := m[key]
	if exists {
		if v, ok := value.(T); ok {
			return v
		}
	}

	var result any
	switch v := value.(type) {
	case float64:
		if _, ok := any(zero).(int64); ok {
			result = int64(v)
		} else if _, ok := any(zero).(bool); ok {
			result = v != 0
		} else {
			return zero
		}
	case int64:
		if _, ok := any(zero).(float64); ok {
			result = float64(v)
		} else if _, ok := any(zero).(bool); ok {
			result = v != 0
		} else {
			return zero
		}
	case string:
		if _, ok := any(zero).(bool); ok {
			lowerVal := strings.ToLower(v)
			if lowerVal == "true" || lowerVal == "1" {
				result = true
			} else if lowerVal == "false" || lowerVal == "0" {
				result = false
			} else {
				return zero
			}
		} else {
			result = v
		}
	case bool:
		if _, ok := any(zero).(string); ok {
			result = fmt.Sprintf("%v", v) // 转成 "true" / "false"
		} else {
			result = v
		}
	default:
		return zero
	}

	if finalValue, ok := result.(T); ok {
		return finalValue
	}

	return zero
}

// GetMapValue 获取map的值
func GetMapValue[T any](m map[string]interface{}, key string) T {
	var zero T

	value, exists := m[key]
	if exists {
		v := reflect.ValueOf(value)
		if v.Type().ConvertibleTo(reflect.TypeOf(zero)) {
			return v.Convert(reflect.TypeOf(zero)).Interface().(T)
		}
	}

	return zero
}

func ConvertToRestfulURL(url string) string {
	re := regexp.MustCompile(`(^.+?/[^/]+)/\d+$`)
	return re.ReplaceAllString(url, `$1/:id`)
}
func GetRequestPath(path, prefix string) (uri string, id int64) {
	uri = strings.TrimPrefix(path, prefix)
	re := regexp.MustCompile(`^(.*)/(\d+)$`)
	matches := re.FindStringSubmatch(uri)
	if len(matches) == 3 {
		uri = matches[1]
		id = StringToInt64(matches[2])
	}

	return
}

// GetPlatform 获取操作系统
func GetPlatform(userAgent string) string {
	ua := strings.ToLower(userAgent)

	// 移动端
	if strings.Contains(ua, "android") {
		return "Android"
	} else if strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad") || strings.Contains(ua, "ipod") {
		return "iOS"
	}

	// 桌面端
	if strings.Contains(ua, "windows") {
		return "Windows"
	} else if strings.Contains(ua, "macintosh") || strings.Contains(ua, "mac os") {
		return "MacOS"
	} else if strings.Contains(ua, "linux") {
		return "Linux"
	}

	return "Unknown"
}

// GetBrowser 获取浏览器类型
func GetBrowser(userAgent string) string {
	ua := strings.ToLower(userAgent)

	if strings.Contains(ua, "chrome") && !strings.Contains(ua, "edg") {
		return "Google Chrome"
	} else if strings.Contains(ua, "edg") {
		return "Microsoft Edge"
	} else if strings.Contains(ua, "firefox") {
		return "Mozilla Firefox"
	} else if strings.Contains(ua, "safari") && !strings.Contains(ua, "chrome") {
		return "Apple Safari"
	} else if strings.Contains(ua, "opr") || strings.Contains(ua, "opera") {
		return "Opera"
	} else if strings.Contains(ua, "msie") || strings.Contains(ua, "trident") {
		return "Internet Explorer"
	}

	return "Unknown"
}
