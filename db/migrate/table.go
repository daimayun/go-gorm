package migrate

import "github.com/daimayun/go-gorm/db"

// GetTableList 获取数据库表列表
func GetTableList() ([]string, error) {
	return db.DB.Migrator().GetTables()
}

// HasTable 表是否存在 [表名 string | model struct ptr]
func HasTable(dst interface{}) bool {
	return db.DB.Migrator().HasTable(dst)
}

// DropTable 删除数据库表
func DropTable(dst ...interface{}) error {
	return db.DB.Migrator().DropTable(dst...)
}

// RenameTable 重命名表名
func RenameTable(oldName, newName interface{}) error {
	return db.DB.Migrator().RenameTable(oldName, newName)
}
