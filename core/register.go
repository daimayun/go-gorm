package core

import (
	"fmt"
	"reflect"

	"go-gorm/core/table"
	"go-gorm/core/tool"
)

// RegisterModel 注册模型
func RegisterModel(models ...interface{}) (err error) {
	for _, model := range models {
		value := reflect.ValueOf(model)
		_type := reflect.Indirect(value).Type()
		if value.Kind() != reflect.Ptr {
			err = fmt.Errorf("cannot use non-ptr model struct `%s`", tool.GetStructFullName(_type))
			return
		}
		if _type.Kind() == reflect.Ptr {
			err = fmt.Errorf("only allow ptr model struct, it looks you use two reference to the struct `%s`", _type)
			return
		}
		if value.Elem().Kind() == reflect.Slice {
			value = reflect.New(value.Elem().Type().Elem())
		}
		if options := table.GetOptions(value); options != "" {
			if err = DB.Set("gorm:table_options", options).AutoMigrate(model); err != nil {
				panic(err)
			}
			return
		}
		if err = DB.AutoMigrate(model); err != nil {
			panic(err)
		}
		return
	}
	return
}
