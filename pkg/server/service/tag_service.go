package service

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/narcissus1949/narcissus-blog/internal/database/mysql"
	cerr "github.com/narcissus1949/narcissus-blog/internal/error"
	"github.com/narcissus1949/narcissus-blog/internal/model"
	"github.com/narcissus1949/narcissus-blog/pkg/dto"
	"github.com/narcissus1949/narcissus-blog/pkg/server/dao"
	"github.com/narcissus1949/narcissus-blog/pkg/vo"
	"go.uber.org/zap"
)

var TagService = new(tagService)

type tagService struct {
}

// ListAllTag 获取所有标签
func (s *tagService) ListAllTag(ctx *gin.Context) ([]vo.TagVo, error) {
	tagList, err := dao.TagDao.ListAllTag(ctx)
	if err != nil {
		zap.L().Error("Failed to list tag name", zap.Error(err))
		return nil, err
	}
	var tagVoList []vo.TagVo
	for _, tag := range tagList {
		tagVoList = append(tagVoList, vo.TagVo{
			ID:          tag.ID,
			Name:        tag.Name,
			CreatedTime: tag.CreatedTime.Format("2006-01-02 15:04:05"),
			UpdatedTime: tag.UpdatedTime.Format("2006-01-02 15:04:05"),
		})
	}
	return tagVoList, nil
}

// ListTag 获取标签列表，分页、条件
func (s *tagService) ListTag(c *gin.Context, tagListDto dto.TagListDto) (*vo.TagListVo, error) {
	// 分页查询 - 获取标签列表
	var tagList []model.ArticleTag
	var tagTotal int64
	txErr := mysql.RunDBTransaction(c, func() error {
		var listErr error
		tagList, listErr = dao.TagDao.ListTag(c, tagListDto)
		if listErr != nil {
			zap.L().Error("Failed to list tag", zap.Error(listErr))
			return listErr
		}
		// 获取标签总数
		var countErr error
		tagTotal, countErr = dao.TagDao.CountTag(c, tagListDto)
		if countErr != nil {
			zap.L().Error("Failed to count tag", zap.Error(countErr))
			return countErr
		}
		return nil
	})
	if txErr != nil {
		zap.L().Error("Failed to list tag", zap.Error(txErr))
		return nil, txErr
	}

	var tagVoList []vo.TagVo
	for _, tag := range tagList {
		tagVoList = append(tagVoList, vo.TagVo{
			ID:          tag.ID,
			Name:        tag.Name,
			CreatedTime: tag.CreatedTime.Format("2006-01-02 15:04:05"),
			UpdatedTime: tag.UpdatedTime.Format("2006-01-02 15:04:05"),
		})
	}

	pageCount := tagTotal / int64(tagListDto.Pageinate.PageSize)
	if tagTotal%int64(tagListDto.Pageinate.PageSize) != 0 {
		pageCount++
	}
	result := vo.TagListVo{
		TagList: tagVoList,
		Pageinate: dto.Pageinate{
			PageSize:  tagListDto.Pageinate.PageSize,
			PageNum:   tagListDto.Pageinate.PageNum,
			Total:     tagTotal,
			PageCount: int(pageCount),
		},
	}
	return &result, nil
}

// ListTagIdByNameArr 根据标签名列表查询标签ID列表
func (s *tagService) ListTagIdByNameArr(ctx *gin.Context, tagDto dto.TagDto) ([]int64, error) {
	idList, err := dao.TagDao.ListTagIdByNameArr(ctx, tagDto.NameList)
	if err != nil {
		zap.L().Error("Failed to list tag id by name arr", zap.Error(err))
		return idList, err
	}
	if len(idList) != len(tagDto.NameList) {
		zap.L().Error("Tag id count does not equal to tag name count",
			zap.Error(err),
			zap.Strings("nameList", tagDto.NameList),
			zap.Int64s("idList", idList))
		return nil, cerr.New(cerr.ERROR_ARTICLE_TAG_NOT_EXIST)
	}
	return idList, nil
}

// GetTagDetail 获取标签详情
func (s *tagService) GetTagDetail(ctx *gin.Context, tagQueryDto dto.TagQueryDto) (*vo.TagVo, error) {
	var tag *model.ArticleTag
	if tagQueryDto.ID != nil {
		var getErr error
		tag, getErr = dao.TagDao.QueryTagByID(ctx, *tagQueryDto.ID)
		if getErr != nil {
			zap.L().Error("Failed to get tag detail", zap.Error(getErr))
			return nil, getErr
		}
	} else if tagQueryDto.Name != nil {
		var getErr error
		tag, getErr = dao.TagDao.QueryTagByName(ctx, *tagQueryDto.Name)
		if getErr != nil {
			zap.L().Error("Failed to get tag detail", zap.Error(getErr))
			return nil, getErr
		}
	} else {
		return nil, cerr.NewParamError()
	}
	if tag == nil {
		return nil, cerr.New(cerr.ERROR_ARTICLE_TAG_NOT_EXIST)
	}

	return &vo.TagVo{
		ID:          tag.ID,
		Name:        tag.Name,
		CreatedTime: tag.CreatedTime.Format("2006-01-02 15:04:05"),
		UpdatedTime: tag.UpdatedTime.Format("2006-01-02 15:04:05"),
	}, nil
}

// CreateTagList 创建标签 - 批量
func (s *tagService) CreateTagList(ctx *gin.Context, tagDto dto.TagDto) error {
	now := time.Now()
	var tagModels []model.ArticleTag
	for _, name := range tagDto.NameList {
		tagModels = append(tagModels, model.ArticleTag{
			Name:        name,
			CreatedTime: now,
			UpdatedTime: now,
		})
	}
	if err := dao.TagDao.InsertTagBatch(ctx, tagModels); err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return cerr.New(cerr.ERROR_ARTICLE_TAG_EXIST)
		}
		zap.L().Error("Failed to insert tag", zap.Error(err))
		return err
	}
	return nil
}

// UpdateTag 更新标签
func (s *tagService) UpdateTag(ctx *gin.Context, updateDto dto.TagUpdateDto) error {
	txErr := mysql.RunDBTransaction(ctx, func() error {
		tagVo, getErr := s.GetTagDetail(ctx, dto.TagQueryDto{
			ID: &updateDto.ID,
		})
		if getErr != nil {
			return getErr
		}

		if tagVo.Name == updateDto.NewName {
			return cerr.NewParamError("tag name not change")
		}

		var tagModel model.ArticleTag
		tagModel.ID = updateDto.ID
		tagModel.Name = updateDto.NewName
		tagModel.UpdatedTime = time.Now()
		if err := dao.TagDao.UpdateTagByID(ctx, tagModel); err != nil {
			zap.L().Error("Failed to update tag", zap.Error(err), zap.Int64("tag id", updateDto.ID))
			return err
		}
		return nil
	})
	if txErr != nil {
		zap.L().Error("Failed to update tag", zap.Error(txErr))
		return txErr
	}

	return nil
}

func (s *tagService) DeleteTagList(deleteDto dto.TagDto) error {
	if err := dao.TagDao.DeleteTagByNameList(deleteDto.NameList); err != nil {
		zap.L().Error("Failed to delete tag", zap.Error(err), zap.Strings("tag", deleteDto.NameList))
		return err
	}
	return nil
}
