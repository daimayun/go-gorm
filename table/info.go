package table

import (
	"reflect"

	"go-gorm/core/tool"
)

// getEngine Get model struct table engine
func getEngine(val reflect.Value) (ok bool, str string) {
	if ok, val = tool.GetFuncReturnValue(val, "TableEngine"); ok {
		return ok, val.String()
	}
	return
}

// getComment Get model struct table comment
func getComment(val reflect.Value) (ok bool, str string) {
	if ok, val = tool.GetFuncReturnValue(val, "TableComment"); ok {
		return ok, val.String()
	}
	return
}

// getCharset Get model struct table charset
func getCharset(val reflect.Value) (ok bool, str string) {
	if ok, val = tool.GetFuncReturnValue(val, "TableCharset"); ok {
		return ok, val.String()
	}
	return
}

// getCollate Get model struct table collate
func getCollate(val reflect.Value) (ok bool, str string) {
	if ok, val = tool.GetFuncReturnValue(val, "TableCollate"); ok {
		return ok, val.String()
	}
	return
}

// getAutoIncrement Get model struct table auto increment
func getAutoIncrement(val reflect.Value) (ok bool, value int64) {
	if ok, val = tool.GetFuncReturnValue(val, "TableAutoIncrement"); ok {
		return ok, val.Int()
	}
	return
}
