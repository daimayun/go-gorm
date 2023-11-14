package migrate

import (
	"gorm.io/gorm"

	"github.com/daimayun/go-gorm/db"
)

// GetTableList 获取数据库表列表
func GetTableList() ([]string, error) {
	return db.DB.Migrator().GetTables()
}

// CreateTable 创建表
func CreateTable(dst ...interface{}) error {
	return db.DB.Migrator().CreateTable(dst...)
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

func TableType(dst interface{}) (gorm.TableType, error) {
	return db.DB.Migrator().TableType(dst)
}

type tableOption struct {
	name          string
	engine        string
	autoIncrement uint64
	charset       string
	collate       string
	comment       string
}

type TableOption func(*tableOption)

func WithTableName(tableName string) TableOption {
	return func(option *tableOption) {
		option.name = tableName
	}
}

func WithEngine(engine string) TableOption {
	return func(option *tableOption) {
		option.engine = engine
	}
}

func WithAutoIncrement(autoIncrement uint64) TableOption {
	return func(option *tableOption) {
		option.autoIncrement = autoIncrement
	}
}

func WithCharset(charset string) TableOption {
	return func(option *tableOption) {
		option.charset = charset
	}
}

func WithCollate(collate string) TableOption {
	return func(option *tableOption) {
		option.collate = collate
	}
}

func WithComment(comment string) TableOption {
	return func(option *tableOption) {
		option.comment = comment
	}
}

// ModifyTable 修改表信息
func ModifyTable(dst interface{}, options ...TableOption) (err error) {
	if len(options) == 0 {
		return
	}
	opt := &tableOption{}
	for _, option := range options {
		option(opt)
	}
	return
}
