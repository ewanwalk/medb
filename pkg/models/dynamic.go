package models

import (
	"fmt"
	"github.com/Ewan-Walker/gorm"
)

func Dyanmic(model interface{}, request map[string]string) func(*gorm.DB) *gorm.DB {

	if len(request) == 0 {
		return func(db *gorm.DB) *gorm.DB {
			return db
		}
	}

	cols := Columns(model)

	return func(db *gorm.DB) *gorm.DB {

		if selects, ok := request["select"]; ok {
			db = db.Select(ValidColumns(selects, model))
		}

		if limit, ok := request["limit"]; ok {
			db = db.Limit(limit)
		}

		if offset, ok := request["offset"]; ok {
			db = db.Offset(offset)
		}

		for field, value := range request {

			colType, ok := cols[field]
			if !ok {
				continue
			}

			switch colType {
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
