package model

const TableNameArticleContent = "article_content"

// ArticleContent mapped from table <article_content>
type ArticleContent struct {
	ID        int    `gorm:"column:id;" json:"id"`
	ArticleID int64  `gorm:"column:article_id;" json:"article_id"`
	Content   string `gorm:"column:content;" json:"content"`
}

// TableName ArticleContent's table name
func (*ArticleContent) TableName() string {
	return TableNameArticleContent
}
