package core

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

var DB *gorm.DB

type Core struct {
	sources  [][]DsnOption
	replicas [][]DsnOption
	pool     []PoolOption
	logger   logger.Interface
	debug    bool
}

// New 初始化
func New(options ...DsnOption) *Core {
	core := &Core{}
	if len(options) > 0 {
		core.sources = append(core.sources, options)
	}
	return core
}

// SetSource 设置主库/写操作
func (c *Core) SetSource(options ...DsnOption) *Core {
	c.sources = append(c.sources, options)
	return c
}

// SetReplica 设置从库/读操作
func (c *Core) SetReplica(options ...DsnOption) *Core {
	c.replicas = append(c.replicas, options)
	return c
}

func (c *Core) SetPoolConfig(options ...PoolOption) *Core {
	c.pool = options
	return c
}

func (c *Core) SetLogger(logger logger.Interface) *Core {
	c.logger = logger
	return c
}

func (c *Core) SetDebug(debug bool) *Core {
	c.debug = debug
	return c
}

func (c *Core) DB() (*gorm.DB, error) {
	sourcesNum := len(c.sources)
	replicasNum := len(c.replicas)
	var sources []string
	if sourcesNum == 0 {
		dsn, _, _ := getDsn()
		sources = append(sources, dsn)
	}
	if sourcesNum > 0 {
		sourcesMap := make(map[string]struct{})
		for _, v := range c.sources {
			dsn, label, _ := getDsn(v...)
			if _, ok := sourcesMap[label]; !ok {
				sourcesMap[label] = struct{}{}
				sources = append(sources, dsn)
			}
		}
		clear(sourcesMap)
	}

	type replica struct {
		dsn  string
		data []interface{}
	}
	replicasMap := make(map[string]replica)
	var replicas []string
	if replicasNum > 0 {
		replicasNum = 0
		for _, v := range c.replicas {
			dsn, label, data := getDsn(v...)
			if _, ok := replicasMap[label]; ok {
				replicasMap[label] = replica{dsn: dsn, data: append(replicasMap[label].data, data...)}
			} else {
				replicasNum++
				replicasMap[label] = replica{dsn: dsn, data: data}
			}
		}
		for k, v := range replicasMap {
			if len(v.data) == 0 {
				replicas = append(replicas, v.dsn)
				delete(replicasMap, k)
			}
		}
	}
	var err error

	DB, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       sources[0], // DSN data source name
		DefaultStringSize:         256,        // string 类型字段的默认长度
		DisableDatetimePrecision:  true,       // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,       // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,       // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,      // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{Logger: c.logger})
	if err != nil {
		return nil, err
	}

	sources = sources[1:]

	// 存在读写分离
	if len(sources) > 0 || replicasNum > 0 {

		dbResolver := dbresolver.Register(dbresolver.Config{TraceResolverMode: true})

		if len(sources) > 0 || len(replicas) > 0 {
			dbResolverConfig := dbresolver.Config{TraceResolverMode: true, Policy: dbresolver.RandomPolicy{}}
			if len(sources) > 0 {
				for _, v := range sources {
					dbResolverConfig.Sources = append(dbResolverConfig.Sources, mysql.Open(v))
				}
			}
			if len(replicas) > 0 {
				for _, v := range replicas {
					dbResolverConfig.Replicas = append(dbResolverConfig.Replicas, mysql.Open(v))
				}
			}
			var data []interface{}
			if len(replicas) == 0 {
				if len(replicasMap) > 0 {
					for k, v := range replicasMap {
						data = v.data
						dbResolverConfig.Replicas = append(dbResolverConfig.Replicas, mysql.Open(v.dsn))
						delete(replicasMap, k)
						break
					}
				}
			}
			dbResolver = dbresolver.Register(dbResolverConfig, data...)
		}

		if len(replicasMap) > 0 {
			for _, v := range replicasMap {
				dbResolverConfig := dbresolver.Config{
					TraceResolverMode: true,
					Replicas:          []gorm.Dialector{mysql.Open(v.dsn)},
				}
				dbResolver.Register(dbResolverConfig, v.data...)
			}
		}

		if err = DB.Use(dbResolver); err != nil {
			return nil, err
		}
	}

	if c.debug {
		DB = DB.Debug()
	}

	return DB, err
}

type poolOptions struct {
	maxIdleConns    int
	maxOpenConns    int
	connMaxLifetime time.Duration
	connMaxIdleTime time.Duration
}

type PoolOption func(options *poolOptions)

func SetMaxIdleConns(num int) PoolOption {
	return func(options *poolOptions) {
		options.maxIdleConns = num
	}
}

func SetMaxOpenConns(num int) PoolOption {
	return func(options *poolOptions) {
		options.maxOpenConns = num
	}
}

func SetConnMaxLifetime(t time.Duration) PoolOption {
	return func(options *poolOptions) {
		options.connMaxLifetime = t
	}
}

func SetConnMaxIdleTime(t time.Duration) PoolOption {
	return func(options *poolOptions) {
		options.connMaxIdleTime = t
	}
}

type dsnOptions struct {
	network   string
	host      string
	port      int
	username  string
	password  string
	database  string
	charset   string
	parseTime string
	loc       string
	data      []interface{}
}

func getDsn(options ...DsnOption) (dsn, label string, data []interface{}) {
	// dsn 默认值
	dsnOption := &dsnOptions{
		host:      "localhost",
		port:      3306,
		username:  "root",
		password:  "root",
		database:  "gorm",
		charset:   "utf8mb4",
		parseTime: "True",
		loc:       "Local",
		network:   "tcp",
	}

	// 设置值
	for _, option := range options {
		option(dsnOption)
	}

	data = dsnOption.data

	dsn = fmt.Sprintf(
		"%s:%s@%s(%s:%d)/%s?charset=%s&parseTime=%s&loc=%s",
		dsnOption.username,
		dsnOption.password,
		dsnOption.network,
		dsnOption.host,
		dsnOption.port,
		dsnOption.database,
		dsnOption.charset,
		dsnOption.parseTime,
		dsnOption.loc,
	)

	label = fmt.Sprintf(
		"%s%d%s%s%s",
		dsnOption.host,
		dsnOption.port,
		dsnOption.username,
		dsnOption.password,
		dsnOption.database,
	)

	return
}

type DsnOption func(*dsnOptions)

func WithHost(host string) DsnOption {
	return func(options *dsnOptions) {
		options.host = host
	}
}

func WithPort(port int) DsnOption {
	return func(options *dsnOptions) {
		options.port = port
	}
}

func WithUsername(username string) DsnOption {
	return func(options *dsnOptions) {
		options.username = username
	}
}

func WithPassword(password string) DsnOption {
	return func(options *dsnOptions) {
		options.password = password
	}
}

func WithDatabase(database string) DsnOption {
	return func(options *dsnOptions) {
		options.database = database
	}
}

func WithCharset(charset string) DsnOption {
	return func(options *dsnOptions) {
		options.charset = charset
	}
}

func WithParseTime(parseTime string) DsnOption {
	return func(options *dsnOptions) {
		options.parseTime = parseTime
	}
}

func WithLoc(loc string) DsnOption {
	return func(options *dsnOptions) {
		options.loc = loc
	}
}

func WithNetwork(network string) DsnOption {
	return func(options *dsnOptions) {
		options.network = network
	}
}

func WithData(data interface{}) DsnOption {
	return func(options *dsnOptions) {
		options.data = append(options.data, data)
	}
}
