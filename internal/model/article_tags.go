package model

import (
	"time"
)

const TableNameArticleTag = "article_tags"

// ArticleTag mapped from table <article_tags>
type ArticleTag struct {
	ID          int64     `gorm:"column:id;" json:"id"`
	Name        string    `gorm:"column:name;" json:"name"`
	CreatedTime time.Time `gorm:"column:created_time;" json:"created_time"`
	UpdatedTime time.Time `gorm:"column:updated_time;" json:"updated_time"`
}

// TableName ArticleTag's table name
func (*ArticleTag) TableName() string {
	return TableNameArticleTag
}
