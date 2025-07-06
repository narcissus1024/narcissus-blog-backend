package dao

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/narcissus1949/narcissus-blog/internal/database/mysql"
	"github.com/narcissus1949/narcissus-blog/internal/model"
	"github.com/narcissus1949/narcissus-blog/pkg/dto"
)

var CategoryDao = &categoryrDao{}

type categoryrDao struct {
}

func (d *categoryrDao) ListAllCategory(ctx *gin.Context) ([]model.ArticleCategory, error) {
	var categoryList []model.ArticleCategory
	res := mysql.GetDBFromContext(ctx).Table(model.TableNameArticleCategory).Find(&categoryList)
	return categoryList, res.Error
}

func (d *categoryrDao) ListCategory(ctx *gin.Context, req dto.CategoryListDto) ([]model.ArticleCategory, error) {
	var categoryList []model.ArticleCategory
	db := mysql.GetDBFromContext(ctx).Table(model.TableNameArticleCategory)
	if len(req.NameListFormat) > 0 {
		db = db.Where("name in ?", req.NameListFormat)
	}
	res := db.Order("name").
		Scopes(dto.Paginate(req.Pageinate)).
		Find(&categoryList)
	return categoryList, res.Error
}

func (d *categoryrDao) CountCategory(ctx *gin.Context, req dto.CategoryListDto) (int64, error) {
	var count int64
	db := mysql.GetDBFromContext(ctx).Table(model.TableNameArticleCategory)
	if len(req.NameListFormat) > 0 {
		db = db.Where("name in ?", req.NameListFormat)
	}
	res := db.Count(&count)
	return count, res.Error
}

func (d *categoryrDao) QueryCategoryByName(ctx *gin.Context, categoryName string) (*model.ArticleCategory, error) {
	var category model.ArticleCategory
	res := mysql.GetDBFromContext(ctx).Table(model.TableNameArticleCategory).Where("name = ?", categoryName).First(&category)
	if res.RowsAffected == 0 {
		return nil, nil
	}
	return &category, res.Error
}

func (d *categoryrDao) QueryCategoryByID(ctx *gin.Context, id int64) (*model.ArticleCategory, error) {
	var category model.ArticleCategory
	res := mysql.GetDBFromContext(ctx).Table(model.TableNameArticleCategory).Where("id = ?", id).First(&category)
	if res.RowsAffected == 0 {
		return nil, nil
	}
	return &category, res.Error
}

func (d *categoryrDao) QueryCategoryIDByName(ctx *gin.Context, categoryName string) (int, error) {
	if len(strings.TrimSpace(categoryName)) == 0 {
		return -1, errors.New("category name invalide")
	}
	var id int
	res := mysql.GetDBFromContext(ctx).Table(model.TableNameArticleCategory).Select("id").Where("name = ?", categoryName).Find(&id)
	if res.RowsAffected == 0 {
		return -1, fmt.Errorf("category %s does not exist", categoryName)
	}
	return id, res.Error
}

func (d *categoryrDao) InsertCategoryBatch(ctx *gin.Context, categoryList []model.ArticleCategory) error {
	res := mysql.GetDBFromContext(ctx).Table(model.TableNameArticleCategory).CreateInBatches(categoryList, 100)
	return res.Error
}

func (d *categoryrDao) UpdateCategoryByID(ctx *gin.Context, category model.ArticleCategory) error {
	res := mysql.GetDBFromContext(ctx).Table(model.TableNameArticleCategory).
		Select("name", "updated_time").
		Where("id = ?", category.ID).
		Updates(category)
	return res.Error
}

func (d *categoryrDao) DeleteCategoryByNameList(ctx *gin.Context, categoryNameList []string) error {
	if len(categoryNameList) == 0 {
		return errors.New("ids is empty")
	}
	res := mysql.GetDBFromContext(ctx).Table(model.TableNameArticleCategory).Where("name in ?", categoryNameList).Delete(&model.ArticleCategory{})
	return res.Error
}
