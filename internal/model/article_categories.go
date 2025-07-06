package model

import (
	"time"
)

const TableNameArticleCategory = "article_categories"

// ArticleCategory mapped from table <article_categories>
type ArticleCategory struct {
	ID          int64     `gorm:"column:id;" json:"id"`
	Name        string    `gorm:"column:name;" json:"name"`
	CreatedTime time.Time `gorm:"column:created_time;" json:"created_time"`
	UpdatedTime time.Time `gorm:"column:updated_time;" json:"updated_time"`
}

// TableName ArticleCategory's table name
func (*ArticleCategory) TableName() string {
	return TableNameArticleCategory
}
