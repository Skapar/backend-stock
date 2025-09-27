package database

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

func ReplaceReturningWildcard(cols []string, sql string) string {
	colsString := strings.Join(cols, ", ")
	re := regexp.MustCompile(`(?i)returning \*`)

	return re.ReplaceAllString(sql, fmt.Sprintf("RETURNING %s", colsString))
}

func ReplaceSelectWildcard(cols []string, sql string) string {
	colsString := strings.Join(cols, ", ")
	re := regexp.MustCompile(`(?i)select \*`)

	return re.ReplaceAllString(sql, fmt.Sprintf("SELECT %s", colsString))
}

func GetEntityDBFilds[T any](entity T) []string {
	var t reflect.Type
	switch reflect.TypeOf(entity).Kind() {
	case reflect.Pointer, reflect.Struct:
		if reflect.TypeOf(entity).Elem().Kind() == reflect.Slice {
			t = reflect.TypeOf(entity).Elem().Elem().Elem()
			break
		}
		t = reflect.TypeOf(entity).Elem()
	case reflect.Slice, reflect.Array:
		t = reflect.TypeOf(entity).Elem().Elem()
	}

	cols := make([]string, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("col")
		if tag != "-" && tag != "" {
			cols = append(cols, tag)
		}
	}
	return cols
}
