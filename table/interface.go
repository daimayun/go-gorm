package table

// ITableComment 表备注
type ITableComment interface {
	TableComment() string
}

// ITableEngine 表引擎
type ITableEngine interface {
	TableEngine() string
}

// ITableCharset 表字符集
type ITableCharset interface {
	TableCharset() string
}

// ITableCollate 表排序规则
type ITableCollate interface {
	TableCollate() string
}

// ITableAutoIncrement 表主键自增长起始值
type ITableAutoIncrement interface {
	TableAutoIncrement() int64
}
