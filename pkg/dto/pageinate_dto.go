package dto

import (
	"gorm.io/gorm"
)

type Pageinate struct {
	PageSize  int   `json:"page_size" form:"page_size" default:"20" binding:"gte=10,lte=50"` // 分页大小
	PageNum   int   `json:"page_num" form:"page_num" default:"1" binding:"gte=1"`            // 当前页
	Total     int64 `json:"total" form:"total"`                                              // 总数据量
	PageCount int   `json:"page_count" form:"page_count"`                                    // 总页数
}

func Paginate(p Pageinate) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page := p.PageNum
		if page <= 0 {
			page = 1
		}

		pageSize := p.PageSize
		switch {
		case pageSize > 50:
			pageSize = 50
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
