package models

import (
	"fmt"
	"github.com/ewanwalk/gorm"
)

func Dynamic(model interface{}, request map[string]string) func(*gorm.DB) *gorm.DB {

	if len(request) == 0 {
		return func(db *gorm.DB) *gorm.DB {
			return db.Model(model)
		}
	}

	cols := Columns(model)

	return func(db *gorm.DB) *gorm.DB {

		db = db.Model(model)

		// select
		if selects, ok := request["select"]; ok {
			db = db.Select(ValidColumns(selects, model))
		}

		// limit
		if limit, ok := request["limit"]; ok {
			db = db.Limit(limit)
		}

		// offset
		if offset, ok := request["offset"]; ok {
			db = db.Offset(offset)
		}

		// order by
		if order, ok := request["order"]; ok {
			method := "asc"
			if m, ok := request["method"]; ok {
				method = m
			}

			db = db.Order(fmt.Sprintf(`%s %s`, order, method))
		}

		// TODO be able to handle date ranges (?)

		for field, value := range request {

			colType, ok := cols[field]
			if !ok {
				continue
			}

			switch colType {
			case "Time":
			case "*time.Time":
			case "time.Time":
				continue
			case "string":
				db = db.Where(fmt.Sprintf("%s LIKE ?", field), fmt.Sprintf("%%%s%%", value))
			case "int":
			case "int64":
				db = db.Where(fmt.Sprintf("%s = ?", field), value)
			}

		}

		return db
	}
}

// DynamicTotal
// used when want we want to obtain a count of every single item available
func DynamicTotal(model interface{}, request map[string]string) func(*gorm.DB) *gorm.DB {

	delete(request, "limit")
	delete(request, "offset")

	return Dynamic(model, request)
}
