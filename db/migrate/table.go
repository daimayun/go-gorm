package migrate

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"

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
		if tableName == "" {
			err = errors.New("表名不能为空")
			return
		}
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

type TableInfoData struct {
	Name          string     `json:"Name"`
	Engine        string     `json:"Engine"`
	Version       string     `json:"Version"`
	RowFormat     string     `json:"Row_format"`
	Rows          int64      `json:"Rows"`
	AvgRowLength  int64      `json:"Avg_row_length"`
	DataLength    int64      `json:"Data_length"`
	MaxDataLength int64      `json:"Max_data_length"`
	IndexLength   int64      `json:"Index_length"`
	DataFree      int64      `json:"Data_free"`
	AutoIncrement int64      `json:"Auto_increment"`
	CreateTime    *time.Time `json:"Create_time"`
	UpdateTime    *time.Time `json:"Update_time"`
	CheckTime     *time.Time `json:"Check_time"`
	Collation     string     `json:"Collation"`
	Checksum      string     `json:"Checksum"`
	CreateOptions string     `json:"Create_options"`
	Comment       string     `json:"Comment"`
}

func TableInfo(dst interface{}) (data TableInfoData, err error) {
	var tableName string
	if tableName, err = getTableName(dst); err != nil {
		return
	}
	var val map[string]interface{}
	if err = db.DB.Raw(fmt.Sprintf("SHOW TABLE STATUS LIKE '%s'", tableName)).Scan(&val).Error; err != nil {
		return
	}
	fmt.Println(val)
	return
}

type tableOption struct {
	tableName     string
	engine        string
	autoIncrement uint64
	charset       string
	collate       string
	comment       string
}

type TableOption func(*tableOption)

func WithNewTableName(tableName string) TableOption {
	return func(option *tableOption) {
		option.tableName = tableName
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

// AlterTable 修改表信息
/*
	ALTER TABLE `abc` RENAME `abc123`, ENGINE = 'MyISAM', AUTO_INCREMENT = 10086, CHARSET = utf8mb4, COLLATE = utf8mb4_general_ci, COMMENT 'abc123表';
*/
func AlterTable(dst interface{}, options ...TableOption) (err error) {
	if len(options) == 0 {
		return
	}

	var tableName string
	if tableName, err = getTableName(dst); err != nil {
		return
	}

	opt := &tableOption{}
	for _, option := range options {
		option(opt)
	}

	exec := false
	sql := fmt.Sprintf("ALTER TABLE `%s`", tableName)
	sign := ""

	if opt.tableName != "" && opt.tableName != tableName {
		sql += fmt.Sprintf("%s RENAME `%s`", sign, opt.tableName)
		sign = ","
		exec = true
	}

	if opt.engine != "" {
		sql += fmt.Sprintf("%s ENGINE = '%s'", sign, opt.engine)
		sign = ","
		exec = true
	}

	if opt.autoIncrement > 0 {
		sql += fmt.Sprintf("%s AUTO_INCREMENT = %d", sign, opt.autoIncrement)
		sign = ","
		exec = true
	}

	if opt.charset != "" {
		sql += fmt.Sprintf("%s CHARSET = %s", sign, opt.charset)
		sign = ","
		exec = true
	}

	if opt.collate != "" {
		sql += fmt.Sprintf("%s COLLATE = %s", sign, opt.collate)
		sign = ","
		exec = true
	}

	if opt.comment != "" {
		sql += fmt.Sprintf("%s COMMENT = '%s'", sign, opt.comment)
		sign = ","
		exec = true
	}

	if !exec {
		return
	}

	if err = db.DB.Exec(sql).Error; err != nil {
		return
	}

	return
}
