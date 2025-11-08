package service

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/narcissus1949/narcissus-blog/internal/database/mysql"
	cerr "github.com/narcissus1949/narcissus-blog/internal/error"
	"github.com/narcissus1949/narcissus-blog/internal/logger"
	"github.com/narcissus1949/narcissus-blog/internal/model"
	"github.com/narcissus1949/narcissus-blog/pkg/dto"
	"github.com/narcissus1949/narcissus-blog/pkg/server/dao"
	"github.com/narcissus1949/narcissus-blog/pkg/vo"
	"go.uber.org/zap"
)

var CategoryService = new(categoryService)

type categoryService struct {
}

// ListAllCategory 获取所有分类
func (s *categoryService) ListAllCategory(ctx *gin.Context) ([]vo.CategoryVo, error) {
	l := logger.FromContext(ctx.Request.Context())
	categoryList, err := dao.CategoryDao.ListAllCategory(ctx)
	if err != nil {
		l.Error("Failed to list all category", zap.Error(err))
		return nil, err
	}
	resp := []vo.CategoryVo{}
	for _, category := range categoryList {
		resp = append(resp, vo.CategoryVo{
			ID:          category.ID,
			Name:        category.Name,
			UpdatedTime: category.UpdatedTime.Format("2006-01-02 15:04:05"),
			CreatedTime: category.CreatedTime.Format("2006-01-02 15:04:05"),
		})
	}
	return resp, nil
}

// ListCategory 获取分类列表 - 分页、条件
func (s *categoryService) ListCategory(ctx *gin.Context, categoryDto dto.CategoryListDto) (*vo.CategoryListVo, error) {
	l := logger.FromContext(ctx.Request.Context())
	// 分页查询 - 获取分类列表
	var categoryList []model.ArticleCategory
	var total int64
	txErr := mysql.RunDBTransaction(ctx, func() error {
		var listErr error
		categoryList, listErr = dao.CategoryDao.ListCategory(ctx, categoryDto)
		if listErr != nil {
			l.Error("Failed to list category", zap.Error(listErr), zap.String("category", categoryDto.NameList))
			return listErr
		}
		// 获取分类总数
		var countErr error
		total, countErr = dao.CategoryDao.CountCategory(ctx, categoryDto)
		if countErr != nil {
			l.Error("Failed to count category", zap.Error(countErr), zap.String("category", categoryDto.NameList))
			return countErr
		}
		return nil
	})
	if txErr != nil {
		l.Error("Failed to run db transaction", zap.Error(txErr))
		return nil, txErr
	}

	voList := []vo.CategoryVo{}
	for _, category := range categoryList {
		voList = append(voList, vo.CategoryVo{
			ID:          category.ID,
			Name:        category.Name,
			UpdatedTime: category.UpdatedTime.Format("2006-01-02 15:04:05"),
			CreatedTime: category.CreatedTime.Format("2006-01-02 15:04:05"),
		})
	}

	pageCount := total / int64(categoryDto.Pageinate.PageSize)
	if total%int64(categoryDto.Pageinate.PageSize) != 0 {
		pageCount++
	}
	result := vo.CategoryListVo{
		CategoryList: voList,
		Pageinate: dto.Pageinate{
			PageSize:  categoryDto.Pageinate.PageSize,
			PageNum:   categoryDto.Pageinate.PageNum,
			Total:     total,
			PageCount: int(pageCount),
		},
	}
	return &result, nil
}

// GetCategoryIDByName 根据分类名获取分类ID
func (s *categoryService) GetCategoryIDByName(ctx *gin.Context, categoryName string) (int, error) {
	l := logger.FromContext(ctx.Request.Context())

	if err := dto.CommonValidateName(categoryName, dto.CATEGORY_MIN_LEN, dto.CATEGORY_MAX_LEN); err != nil {
		return -1, cerr.NewParamError()
	}

	id, err := dao.CategoryDao.QueryCategoryIDByName(ctx, categoryName)
	if err != nil {
		l.Error("Failed to query category id by name", zap.Error(err), zap.String("category", categoryName))
		return -1, err
	}
	if id < 0 {
		return -1, cerr.New(cerr.ERROR_ARTICLE_CATEGORY_NOT_EXIST)
	}
	return id, nil
}

// GetCategoryDetail 获取分类详情
func (s *categoryService) GetCategoryDetail(ctx *gin.Context, categoryQueryDto dto.CategoryQueryDto) (*vo.CategoryVo, error) {
	l := logger.FromContext(ctx.Request.Context())
	var category *model.ArticleCategory
	if categoryQueryDto.ID != nil {
		var getErr error
		category, getErr = dao.CategoryDao.QueryCategoryByID(ctx, *categoryQueryDto.ID)
		if getErr != nil {
			l.Error("Failed to get category detail", zap.Error(getErr))
			return nil, getErr
		}
	} else if categoryQueryDto.Name != nil {
		var getErr error
		category, getErr = dao.CategoryDao.QueryCategoryByName(ctx, *categoryQueryDto.Name)
		if getErr != nil {
			l.Error("Failed to get category detail", zap.Error(getErr))
			return nil, getErr
		}
	} else {
		return nil, cerr.NewParamError()
	}
	if category == nil {
		return nil, cerr.New(cerr.ERROR_ARTICLE_CATEGORY_NOT_EXIST)
	}

	return &vo.CategoryVo{
		ID:          category.ID,
		Name:        category.Name,
		CreatedTime: category.CreatedTime.Format("2006-01-02 15:04:05"),
		UpdatedTime: category.UpdatedTime.Format("2006-01-02 15:04:05"),
	}, nil
}

// CreateCategoryList 创建分类 - 批量
func (s *categoryService) CreateCategoryList(ctx *gin.Context, categoryDto dto.CategoryDto) error {
	l := logger.FromContext(ctx.Request.Context())
	now := time.Now()
	var categoryModels []model.ArticleCategory
	for _, name := range categoryDto.NameList {
		categoryModels = append(categoryModels, model.ArticleCategory{
			Name:        name,
			CreatedTime: now,
			UpdatedTime: now,
		})
	}

	if err := dao.CategoryDao.InsertCategoryBatch(ctx, categoryModels); err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return cerr.New(cerr.ERROR_ARTICLE_CATEGORY_EXIST)
		}
		l.Error("Failed to insert category", zap.Error(err))
		return err
	}
	return nil
}

// UpdateCategory 更新分类
func (s *categoryService) UpdateCategory(ctx *gin.Context, updateDto dto.CategoryUpdateDto) error {
	l := logger.FromContext(ctx.Request.Context())
	txErr := mysql.RunDBTransaction(ctx, func() error {
		categoryVo, getErr := s.GetCategoryDetail(ctx, dto.CategoryQueryDto{
			ID: &updateDto.ID,
		})
		if getErr != nil {
			l.Error("Failed to get category detail", zap.Error(getErr))
			return getErr
		}

		if categoryVo.Name == updateDto.NewName {
			return cerr.NewParamError("category name not change")
		}

		var categoryModel model.ArticleCategory
		categoryModel.ID = updateDto.ID
		categoryModel.Name = updateDto.NewName
		categoryModel.UpdatedTime = time.Now()
		if err := dao.CategoryDao.UpdateCategoryByID(ctx, categoryModel); err != nil {
			l.Error("Failed to update category", zap.Error(err), zap.Int64("category id", updateDto.ID))
			return err
		}
		return nil
	})
	if txErr != nil {
		l.Error("Failed to update category", zap.Error(txErr))
		return txErr
	}

	return nil
}

// DeleteCategory 删除分类 - 批量
func (s *categoryService) DeleteCategoryList(ctx *gin.Context, deleteRequest dto.CategoryDto) error {
	l := logger.FromContext(ctx.Request.Context())
	if err := dao.CategoryDao.DeleteCategoryByNameList(ctx, deleteRequest.NameList); err != nil {
		l.Error("Failed to delete category", zap.Error(err), zap.Strings("category", deleteRequest.NameList))
		return err
	}
	return nil
}
