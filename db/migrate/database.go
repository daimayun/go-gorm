package migrate

import "github.com/daimayun/go-gorm/db"

// GetCurrentDatabase 获取当前的数据库名
func GetCurrentDatabase() string {
	return db.DB.Migrator().CurrentDatabase()
}
