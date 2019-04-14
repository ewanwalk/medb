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
		if tag := m.Field(i).Tag.Get("json"); len(tag) != 0 && tag != "-" {
			if idx := strings.Index(tag, ","); idx > 0 {
				columns[tag[:idx]] = m.Field(i).Type.Name()
			}
		}
	}

	return columns
}
