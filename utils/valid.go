package utils

import (
	"reflect"
	"strconv"
)

type Number any
type String any

func IsValidNumber[T Number](value T) bool {
	switch v := any(value).(type) {
	case int:
		return v > 0
	case int32:
		return v > 0
	case int64:
		return v > 0
	case float64:
		return v > 0
	case float32:
		return v > 0
	case string:
		if num, err := strconv.ParseFloat(v, 64); err == nil {
			return num > 0
		}
	}
	return false
}

func IsValidString[T any](value T) bool {
	switch v := any(value).(type) {
	case string:
		return len(v) > 0
	}

	return false
}

func IsValidPageQuery[T any](value T) bool {
	v := reflect.ValueOf(value)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return false
	}

	page := v.FieldByName("Page")
	if !page.IsValid() {
		return false
	}

	pageSize := v.FieldByName("PageSize")
	if !pageSize.IsValid() {
		return false
	}

	allowedKinds := map[reflect.Kind]bool{
		reflect.Int:    true,
		reflect.Int32:  true,
		reflect.Int64:  true,
		reflect.Uint:   true,
		reflect.Uint32: true,
		reflect.Uint64: true,
	}

	if !allowedKinds[page.Kind()] || !allowedKinds[pageSize.Kind()] {
		return false
	}

	if page.IsZero() || pageSize.IsZero() {
		return false
	}

	return true
}
