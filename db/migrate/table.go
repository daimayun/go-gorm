package migrate

import "github.com/daimayun/go-gorm/db"

// GetTableList 获取表列表
func GetTableList() ([]string, error) {
	return db.DB.Migrator().GetTables()
}
