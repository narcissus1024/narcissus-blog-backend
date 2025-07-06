package dao

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/narcissus1949/narcissus-blog/internal/database/mysql"
	"github.com/narcissus1949/narcissus-blog/internal/model"
	"github.com/narcissus1949/narcissus-blog/internal/utils"
	"github.com/narcissus1949/narcissus-blog/pkg/dto"
	"gorm.io/gorm"
)

var ArticleDao = &articlerDao{}

type articlerDao struct {
}

func (d *articlerDao) InsertArticle(c *gin.Context, article *model.Article) error {
	tx := mysql.GetDBFromContext(c)
	res := tx.Create(article)
	return res.Error
}

func (d *articlerDao) UpdateArticle(c *gin.Context, article *model.Article) (int64, error) {
	if article.ID == 0 {
		return 0, errors.New("article id is required")
	}
	tx := mysql.GetDBFromContext(c)
	res := tx.Model(article).
		Select(
			"title",
			"summary",
			"category_id",
			"author",
			"allow_comment",
			"weight",
			"is_sticky",
			"is_original",
			"original_article_link",
			"status",
			"updated_time").
		Updates(article)
	return res.RowsAffected, res.Error
}

/*
select a.*, c.name as category_name, group_concat(t.name) as tag_name_list from articles as a
left join article_categories as c on a.category_id = c.id
left join article_tags as t on FIND_IN_SET(t.id,a.tags_id) > 0
where (a.type in (0,1) and a.is_original = 1) and ((a.type = 0 and c.name = "golang" and t.name in ("基础") ) or (a.type = 1))
group by a.id
*/
func (d *articlerDao) ListArticle(ctx *gin.Context, articleListRequest dto.ArticleListDto) ([]model.ArticleDetail, error) {
	db := mysql.GetDBFromContext(ctx)
	var articleList []model.ArticleDetail
	db = db.Table(model.TableNameArticle + " a").
		Select("a.*, c.name as category_name, GROUP_CONCAT(t.name) as tag_name_list").
		Joins(fmt.Sprintf("left join %s c on a.category_id = c.id", model.TableNameArticleCategory)).
		Joins(fmt.Sprintf("left join %s r on r.article_id = a.id", model.TableNameArticleTagRelation)).
		Joins(fmt.Sprintf("left join %s t on t.id = r.tag_id", model.TableNameArticleTag))

	// 条件查询
	// 不同文章类型(type)公用字段
	{
		baseCond := db.Session(&gorm.Session{NewDB: true})
		baseCond = baseCond.Where("a.type in ?", articleListRequest.Type)
		if len(strings.TrimSpace(articleListRequest.Title)) > 0 {
			baseCond = baseCond.Where("a.title like ?", "%"+articleListRequest.Title+"%")
		}
		if len(strings.TrimSpace(articleListRequest.Author)) > 0 {
			baseCond = baseCond.Where("a.author like ?", "%"+articleListRequest.Author+"%")
		}
		if articleListRequest.IsOriginal != nil {
			baseCond = baseCond.Where("a.is_original = ?", *articleListRequest.IsOriginal)
		}
		if articleListRequest.Status != nil {
			baseCond = baseCond.Where("a.status = ?", *articleListRequest.Status)
		}
		if articleListRequest.StartTime > 0 {
			baseCond = baseCond.Where("a.created_time >= ?", articleListRequest.StartTime)
		}
		if articleListRequest.EndTime > 0 {
			baseCond = baseCond.Where("a.updated_time <= ?", articleListRequest.EndTime)
		}
		db = db.Where(baseCond)
	}
	// 博文支持分类、标签;随笔不支持分类、标签
	{
		baseCond := db.Session(&gorm.Session{NewDB: true})
		if utils.ArrayExistInt(articleListRequest.Type, utils.ARTICLE_TYPE_POST) {
			postCond := baseCond.Session(&gorm.Session{NewDB: true})
			postCond = postCond.Where("a.type = ?", utils.ARTICLE_TYPE_POST)
			if len(strings.TrimSpace(articleListRequest.Category)) > 0 {
				postCond = postCond.Where("c.name = ?", articleListRequest.Category)
			}
			if len(articleListRequest.Tags) > 0 {
				postCond = postCond.Where("t.name in ?", articleListRequest.Tags)
			}
			baseCond = baseCond.Or(postCond)
		}
		if utils.ArrayExistInt(articleListRequest.Type, utils.ARTICLE_TYPE_ESSAY) {
			essayCond := baseCond.Session(&gorm.Session{NewDB: true})
			essayCond = essayCond.Where("a.type = ?", utils.ARTICLE_TYPE_ESSAY)
			baseCond = baseCond.Or(essayCond)
		}
		db = db.Where(baseCond)
	}

	db = db.Group("a.id")
	// 排序
	db = db.Order("a.created_time desc")
	// 分页查询
	res := db.Scopes(dto.Paginate(articleListRequest.Pageinate)).Find(&articleList)

	return articleList, res.Error
}

// todo 目前只支持分类计数、标签计数和全量计数
func (d *articlerDao) CountArticle(ctx *gin.Context, articleListRequest dto.ArticleListDto) (int64, error) {
	db := mysql.GetDBFromContext(ctx)
	var totalArticle int64
	db = db.Table(model.TableNameArticle + " a").
		Select("a.*, c.name as category_name, GROUP_CONCAT(t.name) as tag_name_list").
		Joins(fmt.Sprintf("left join %s c on a.category_id = c.id", model.TableNameArticleCategory)).
		Joins(fmt.Sprintf("left join %s r on r.article_id = a.id", model.TableNameArticleTagRelation)).
		Joins(fmt.Sprintf("left join %s t on t.id = r.tag_id", model.TableNameArticleTag))

	// 条件查询
	// 不同文章类型(type)公用字段
	{
		baseCond := db.Session(&gorm.Session{NewDB: true})
		baseCond = baseCond.Where("a.type in ?", articleListRequest.Type)
		if len(strings.TrimSpace(articleListRequest.Title)) > 0 {
			baseCond = baseCond.Where("a.title like ?", "%"+articleListRequest.Title+"%")
		}
		if len(strings.TrimSpace(articleListRequest.Author)) > 0 {
			baseCond = baseCond.Where("a.author like ?", "%"+articleListRequest.Author+"%")
		}
		if articleListRequest.IsOriginal != nil {
			baseCond = baseCond.Where("a.is_original = ?", *articleListRequest.IsOriginal)
		}
		if articleListRequest.Status != nil {
			baseCond = baseCond.Where("a.status = ?", *articleListRequest.Status)
		}
		if articleListRequest.StartTime > 0 {
			baseCond = baseCond.Where("a.created_time >= ?", articleListRequest.StartTime)
		}
		if articleListRequest.EndTime > 0 {
			baseCond = baseCond.Where("a.updated_time <= ?", articleListRequest.EndTime)
		}
		db = db.Where(baseCond)
	}
	// 博文支持分类、标签;随笔不支持分类、标签
	{
		baseCond := db.Session(&gorm.Session{NewDB: true})
		if utils.ArrayExistInt(articleListRequest.Type, utils.ARTICLE_TYPE_POST) {
			postCond := baseCond.Session(&gorm.Session{NewDB: true})
			postCond = postCond.Where("a.type = ?", utils.ARTICLE_TYPE_POST)
			if len(strings.TrimSpace(articleListRequest.Category)) > 0 {
				postCond = postCond.Where("c.name = ?", articleListRequest.Category)
			}
			if len(articleListRequest.Tags) > 0 {
				postCond = postCond.Where("t.name in ?", articleListRequest.Tags)
			}
			baseCond = baseCond.Or(postCond)
		}
		if utils.ArrayExistInt(articleListRequest.Type, utils.ARTICLE_TYPE_ESSAY) {
			essayCond := baseCond.Session(&gorm.Session{NewDB: true})
			essayCond = essayCond.Where("a.type = ?", utils.ARTICLE_TYPE_ESSAY)
			baseCond = baseCond.Or(essayCond)
		}
		db = db.Where(baseCond)
	}

	db = db.Group("a.id")
	res := db.Count(&totalArticle)

	return totalArticle, res.Error
}

/*
select a.*,c.name as category_name,ac.content,GROUP_CONCAT(t.name) from articles as a
left join article_categories as c on a.category_id = c.id
left join article_content as ac on a.id = ac.article_id
left join article_tag_relations as r on a.id = r.article_id
left join article_tags as t on r.tag_id = t.id
where a.id = 1;
*/
func (d *articlerDao) QueryArticleDetail(c *gin.Context, id int64) (*model.ArticleDetail, error) {
	if id <= 0 {
		return nil, errors.New("in invalide")
	}
	db := mysql.GetDBFromContext(c)
	var detail model.ArticleDetail
	// 查询文章基本信息和分类名称
	res := db.Table(model.TableNameArticle+" as a").
		Select("a.*, c.name as category_name, ctx.content as content, GROUP_CONCAT(t.name) as tag_name_list").
		Joins(fmt.Sprintf("left join %s c on a.category_id = c.id", model.TableNameArticleCategory)).
		Joins(fmt.Sprintf("left join %s as ctx on a.id = ctx.article_id", model.TableNameArticleContent)).
		Joins(fmt.Sprintf("left join %s r on a.id = r.article_id", model.TableNameArticleTagRelation)).
		Joins(fmt.Sprintf("left join %s t on r.tag_id = t.id", model.TableNameArticleTag)).
		Where("a.id = ?", id).
		Group("a.id").
		First(&detail)

	if res.RowsAffected == 0 {
		return nil, res.Error
	}
	return &detail, nil
}

func (d *articlerDao) DeleteArticleByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return errors.New("ids is empty")
	}
	db := mysql.GetDBFromContext(c)
	return db.Table(model.TableNameArticle).Where("id in ?", ids).Delete(&model.Article{}).Error
}

func (d *articlerDao) DeleteArticleByID(c *gin.Context, id int64) error {
	if id <= 0 {
		return errors.New("id is invalide")
	}
	db := mysql.GetDBFromContext(c)
	return db.Table(model.TableNameArticle).Where("id = ?", id).Delete(&model.Article{}).Error
}
