package migrate

import (
	"errors"
	"fmt"
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

// 获取表名
func getTableName(dst interface{}) (tableName string, err error) {
	v, ok := dst.(string)
	if ok {
		tableName = v
		return
	}
	stmt := &gorm.Statement{DB: db.DB}
	if err = stmt.Parse(dst); err == nil {
		tableName = stmt.Table
		return
	}
	return
}

func TableDDL(dst interface{}) (str string, err error) {
	var tableName string
	if tableName, err = getTableName(dst); err != nil {
		return
	}
	var val map[string]interface{}
	if err = db.DB.Raw(fmt.Sprintf("SHOW CREATE TABLE `%s`", tableName)).Scan(&val).Error; err != nil {
		return
	}
	value, ok := val["Create Table"]
	if ok {
		str = value.(string)
		return
	}
	err = errors.New("table [" + tableName + "] ddl not found")
	return
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
	var tableName string
	// 表名处理
	if v, ok := dst.(string); ok {
		tableName = v
	} else {
		stmt := &gorm.Statement{DB: db.DB}
		if err = stmt.Parse(dst); err == nil {
			tableName = stmt.Table
		} else {
			return err
		}
	}

	// 修改表名
	fmt.Sprintf("ALTER TABLE `%s` RENAME `%s`;", tableName, opt.name)

	return
}
