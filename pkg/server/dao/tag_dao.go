package dao

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/narcissus1949/narcissus-blog/internal/database/mysql"
	"github.com/narcissus1949/narcissus-blog/internal/model"
	"github.com/narcissus1949/narcissus-blog/pkg/dto"
)

var TagDao = &tagDao{}

type tagDao struct {
}

func (d *tagDao) QueryTagByName(ctx *gin.Context, name string) (*model.ArticleTag, error) {
	var tag model.ArticleTag
	res := mysql.GetDBFromContext(ctx).Where("name = ?", name).First(&tag)
	return &tag, res.Error
}

// GetTagDetailByID 根据标签 ID 获取标签详情
func (d *tagDao) QueryTagByID(ctx *gin.Context, id int64) (*model.ArticleTag, error) {
	var tag model.ArticleTag
	res := mysql.GetDBFromContext(ctx).Where("id = ?", id).First(&tag)
	return &tag, res.Error
}

func (d *tagDao) ListAllTag(ctx *gin.Context) ([]model.ArticleTag, error) {
	var tagList []model.ArticleTag
	res := mysql.GetDBFromContext(ctx).Table(model.TableNameArticleTag).Find(&tagList)
	return tagList, res.Error
}

func (d *tagDao) ListTag(c *gin.Context, req dto.TagListDto) ([]model.ArticleTag, error) {
	var tagList []model.ArticleTag
	db := mysql.GetDBFromContext(c)
	db = db.Table(model.TableNameArticleTag)
	if len(req.NameListFormat) > 0 {
		db = db.Where("name in ?", req.NameListFormat)
	}
	res := db.Order("name").
		Scopes(dto.Paginate(req.Pageinate)).
		Find(&tagList)
	return tagList, res.Error
}

// CountTag 统计标签数量
func (d *tagDao) CountTag(c *gin.Context, tagListDto dto.TagListDto) (int64, error) {
	var count int64
	db := mysql.GetDBFromContext(c)
	db = db.Table(model.TableNameArticleTag)
	if len(tagListDto.NameListFormat) > 0 {
		db = db.Where("name in ?", tagListDto.NameListFormat)
	}
	res := db.Count(&count)
	return count, res.Error
}

func (d *tagDao) ListTagIdByNameArr(ctx *gin.Context, nameList []string) ([]int64, error) {
	if len(nameList) == 0 {
		return nil, errors.New("nameList is empty")
	}
	var idList []int64
	res := mysql.GetDBFromContext(ctx).Table(model.TableNameArticleTag).Select("id").Where("name in ?", nameList).Find(&idList)
	return idList, res.Error
}

// GetTagDetail 获取标签详情
func (d *tagDao) GetTagDetail(ctx *gin.Context, nameList []string) (*model.ArticleTag, error) {
	var tag model.ArticleTag
	res := mysql.GetDBFromContext(ctx).Table(model.TableNameArticleTag).Where("name in ?", nameList).First(&tag)
	return &tag, res.Error
}

func (d *tagDao) InsertTagBatch(ctx *gin.Context, tagList []model.ArticleTag) error {
	res := mysql.GetDBFromContext(ctx).Table(model.TableNameArticleTag).CreateInBatches(tagList, 100)
	return res.Error
}

// UpdateTag 根据标签 ID 更新标签信息
func (d *tagDao) UpdateTagByID(ctx *gin.Context, tag model.ArticleTag) error {
	res := mysql.GetDBFromContext(ctx).Table(model.TableNameArticleTag).
		Select("name", "updated_time").
		Where("id = ?", tag.ID).
		Updates(tag)
	return res.Error
}

// DeleteTag 根据标签
func (d *tagDao) DeleteTagByNameList(nameList []string) error {
	if len(nameList) == 0 {
		return errors.New("nameList is empty")
	}
	res := mysql.Client.Table(model.TableNameArticleTag).
		Where("name in ?", nameList).
		Delete(&model.ArticleTag{})
	return res.Error
}
