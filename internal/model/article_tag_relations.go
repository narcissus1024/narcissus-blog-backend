package model

const (
	TableNameArticleTagRelation = "article_tag_relations"
)

// ArticleTagRelation 文章 - 标签关系结构体
type ArticleTagRelation struct {
	ID        int64 `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	ArticleID int64 `json:"article_id" gorm:"column:article_id;not null"`
	TagID     int64 `json:"tag_id" gorm:"column:tag_id;not null"`
}

func (*ArticleTagRelation) TableName() string {
	return TableNameArticleTagRelation
}
