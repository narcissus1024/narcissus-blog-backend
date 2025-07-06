package dao

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/narcissus1949/narcissus-blog/internal/database/mysql"
	"github.com/narcissus1949/narcissus-blog/internal/model"
)

var ArticleContentDao = &articleContentDao{}

type articleContentDao struct {
}

func (d *articleContentDao) InsertContent(c *gin.Context, content *model.ArticleContent) error {
	if content.ArticleID <= 0 {
		return errors.New("article id is invalid")
	}
	tx := mysql.GetDBFromContext(c)
	res := tx.Create(content)
	return res.Error
}

func (d *articleContentDao) UpdateContentByArticleID(c *gin.Context, content *model.ArticleContent) (int64, error) {
	if content == nil || content.ArticleID <= 0 {
		return 0, errors.New("article content params is invalid")
	}
	tx := mysql.GetDBFromContext(c)
	res := tx.Model(content).
		Select("content").
		Where("article_id = ?", content.ArticleID).
		Updates(content)
	return res.RowsAffected, res.Error
}

func (d *articleContentDao) DeleteContentByArticleIDs(c *gin.Context, articleIDs []int64) error {
	if len(articleIDs) <= 0 {
		return errors.New("article id is invalid")
	}
	db := mysql.GetDBFromContext(c)
	return db.Table(model.TableNameArticleContent).Where("article_id in ?", articleIDs).Delete(&model.ArticleContent{}).Error
}
