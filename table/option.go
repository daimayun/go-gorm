package table

import (
	"fmt"
	"reflect"
)

// GetOptions GetTableOptions
func GetOptions(table reflect.Value) (options string) {
	var sign string
	if ok, tableEngine := getEngine(table); ok && tableEngine != "" {
		options += sign + "ENGINE=" + tableEngine
		sign = " "
	}
	if ok, tableAutoIncrement := getAutoIncrement(table); ok && tableAutoIncrement > 0 {
		options += fmt.Sprintf("%sAUTO_INCREMENT=%d", sign, tableAutoIncrement)
		sign = " "
	}
	if ok, tableCharset := getCharset(table); ok && tableCharset != "" {
		options += sign + "CHARSET=" + tableCharset
		sign = " "
	}
	if ok, tableCollate := getCollate(table); ok && tableCollate != "" {
		options += sign + "COLLATE=" + tableCollate
		sign = " "
	}
	if ok, tableComment := getComment(table); ok && tableComment != "" {
		options += sign + "COMMENT='" + tableComment + "'"
		sign = " "
	}
	return
}
