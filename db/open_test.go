package db_test

import (
	"database/sql"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

func low() {
	dsn := "root:root@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	fmt.Println(db)
}

func high() {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       "root:root@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local", // DSN data source name
		DefaultStringSize:         256,                                                                           // string 类型字段的默认长度
		DisableDatetimePrecision:  true,                                                                          // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,                                                                          // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,                                                                          // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,                                                                         // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	fmt.Println(db)
}

func conn() {
	sqlDB, err := sql.Open("mysql", "mydb_dsn")
	if err != nil {
		panic(err)
	}
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	fmt.Println(gormDB)
}

func resolver() {
	db, err := gorm.Open(mysql.Open("db1_dsn"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = db.Use(dbresolver.Register(dbresolver.Config{
		// `db2` 作为 sources，`db3`、`db4` 作为 replicas
		Sources:  []gorm.Dialector{mysql.Open("db2_dsn")},
		Replicas: []gorm.Dialector{mysql.Open("db3_dsn"), mysql.Open("db4_dsn")},
		// sources/replicas 负载均衡策略
		Policy: dbresolver.RandomPolicy{},
	}).Register(dbresolver.Config{
		// `db1` 作为 sources（DB 的默认连接），对于 `User`、`Address` 使用 `db5` 作为 replicas
		Replicas: []gorm.Dialector{mysql.Open("db5_dsn")},
	}, "user", "address").Register(dbresolver.Config{
		// `db6`、`db7` 作为 sources，对于 `orders`、`Product` 使用 `db8` 作为 replicas
		Sources:  []gorm.Dialector{mysql.Open("db6_dsn"), mysql.Open("db7_dsn")},
		Replicas: []gorm.Dialector{mysql.Open("db8_dsn")},
	}, "orders", "product", "secondary"))
	if err != nil {
		panic(err)
	}
}
