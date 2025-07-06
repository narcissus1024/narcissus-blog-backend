package dao

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/narcissus1949/narcissus-blog/internal/database/mysql"
	"github.com/narcissus1949/narcissus-blog/internal/model"
)

var ArticleTagRelationDao = &articleTagRelationDao{}

type articleTagRelationDao struct {
}

// InsertArticleTagRelations 批量插入文章标签关系
func (d *articleTagRelationDao) InsertArticleTagRelations(c *gin.Context, relations []*model.ArticleTagRelation) error {
	if len(relations) == 0 {
		return errors.New("article tag relation is empty")
	}
	tx := mysql.GetDBFromContext(c)
	return tx.CreateInBatches(relations, 100).Error
}

// DeleteArticleTagRelationsBatch 批量删除文章标签关系
func (d *articleTagRelationDao) DeleteArticleTagRelationsByArticleIDs(c *gin.Context, articleIDs []int64) error {
	if len(articleIDs) == 0 {
		return errors.New("article id is empty")
	}
	tx := mysql.GetDBFromContext(c)
	return tx.Where("article_id in ?", articleIDs).Delete(&model.ArticleTagRelation{}).Error
}

// DeleteArticleTagRelationsByArticleNames 批量删除文章标签关系
func (d *articleTagRelationDao) DeleteArticleTagRelations(c *gin.Context, relations []*model.ArticleTagRelation) error {
	if len(relations) == 0 {
		return errors.New("article tag relation is empty")
	}
	tx := mysql.GetDBFromContext(c)

	// 构建批量删除条件
	var conditions []string
	var args []interface{}
	for _, rel := range relations {
		conditions = append(conditions, "(article_id = ? AND tag_id = ?)")
		args = append(args, rel.ArticleID, rel.TagID)
	}

	// 执行批量删除
	res := tx.Where(strings.Join(conditions, " OR "), args...).Delete(&model.ArticleTagRelation{})
	return res.Error
}
