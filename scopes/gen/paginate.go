package gen

import "gorm.io/gen"

// Paginate 分页
func Paginate(page, pageSize int) func(dao gen.Dao) gen.Dao {
	return func(dao gen.Dao) gen.Dao {
		if page == 0 {
			page = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}
		offset := (page - 1) * pageSize
		return dao.Offset(offset).Limit(pageSize)
	}
}
