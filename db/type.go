package db

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// DateTime 时间类型
type DateTime time.Time

// MarshalJSON Json后的数据处理
func (t DateTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(t).In(time.Local).Format(time.DateTime) + `"`), nil
}

// Value 存入数据库
func (t DateTime) Value() (driver.Value, error) {
	return time.Time(t), nil
}

// Slice 切片类型/列表类型
type Slice []any

// Value 存入数据库
func (l Slice) Value() (driver.Value, error) {
	return json.Marshal(l)
}

// Scan 从数据库读出
func (l *Slice) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &l)
}
