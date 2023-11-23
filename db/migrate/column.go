package migrate

import (
	"errors"

	"github.com/daimayun/go-gorm/db"
)

// HasColumn 表字段是否存在
func HasColumn(dst interface{}, field string) bool {
	return db.DB.Migrator().HasColumn(dst, field)
}

func AddColumn(dst interface{}) (err error) {
	return
}

func DeleteColumn(dst interface{}, field string) (err error) {
	if !HasColumn(dst, field) {
		err = errors.New("column [" + field + "] no exist")
		return
	}
	return
}
