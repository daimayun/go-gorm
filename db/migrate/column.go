package migrate

import "github.com/daimayun/go-gorm/db"

// HasColumn 表字段是否存在
func HasColumn(dst interface{}, field string) bool {
	return db.DB.Migrator().HasColumn(dst, field)
}
