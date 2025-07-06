package model

import (
	"time"
)

const TableNameArticle = "articles"

// Article mapped from table <articles>
type Article struct {
	ID                  int64     `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Title               string    `json:"title" gorm:"column:title;not null"`
	Summary             string    `json:"summary" gorm:"column:summary;type:text"`
	Type                uint8     `json:"type" gorm:"column:type;not null"`
	CategoryID          int       `json:"category_id" gorm:"column:category_id;null"`
	Author              string    `json:"author" gorm:"column:author;not null"`
	AllowComment        bool      `json:"allow_comment" gorm:"column:allow_comment"`
	Weight              int       `json:"weight" gorm:"column:weight"`
	IsSticky            bool      `json:"is_sticky" gorm:"column:is_sticky"`
	IsOriginal          bool      `json:"is_original" gorm:"column:is_original"`
	OriginalArticleLink string    `json:"original_article_link" gorm:"column:original_article_link;null"`
	Status              uint8     `json:"status" gorm:"column:status"`
	CreatedTime         time.Time `json:"created_time" gorm:"column:created_time;autoCreateTime"`
	UpdatedTime         time.Time `json:"updated_time" gorm:"column:updated_time;autoUpdateTime"`
}

// TableName Article's table name
func (*Article) TableName() string {
	return TableNameArticle
}

// 文章内容，将分类id和标签id转为名称
type ArticleDetail struct {
	Article
	CategoryName string
	TagNameList  string
	Content      string
}
