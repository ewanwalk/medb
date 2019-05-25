package models

import (
	"reflect"
	"strings"
)

// ValidColumns
// obtain a slice of valid column names based on a string of comma delimited columns provided
// based on the model provided
// Note: this relies on json struct tags: `json:"example"`
func ValidColumns(columns string, model interface{}) []string {

	cols := strings.Split(strings.ToLower(columns), ",")
	exist := make([]string, 0)

	clmns := Columns(model)
	for _, col := range cols {
		if _, ok := clmns[col]; ok {
			exist = append(exist, col)
		}
	}

	return exist
}

// Columns
// obtain all columns and their types
func Columns(model interface{}) map[string]string {

	columns := make(map[string]string, 0)

	m := reflect.ValueOf(model).Type()
	for i := 0; i < m.NumField(); i++ {
		field := m.Field(i)
		if tag := field.Tag.Get("json"); len(tag) != 0 && tag != "-" {

			name := field.Type.Name()
			if field.Type.Kind().String() == "ptr" {
				name = field.Type.Elem().Name()
			}

			if idx := strings.Index(tag, ","); idx > 0 {
				columns[tag[:idx]] = name
			} else {
				columns[tag] = name
			}
		}
	}

	return columns
}
