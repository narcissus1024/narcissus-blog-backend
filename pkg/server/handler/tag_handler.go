package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/narcissus1949/narcissus-blog/pkg/dto"
	"github.com/narcissus1949/narcissus-blog/pkg/server/service"
	resp "github.com/narcissus1949/narcissus-blog/pkg/vo/response"
	"go.uber.org/zap"
)

var TagHandler = new(tagHandler)

type tagHandler struct {
}

func (c *tagHandler) ListAllTag(ctx *gin.Context) {
	tagList, listErr := service.TagService.ListAllTag(ctx)
	if listErr != nil {
		resp.Fail(ctx, listErr)
		return
	}

	resp.OK(ctx, tagList)
}

func (c *tagHandler) ListTag(ctx *gin.Context) {
	var tagListDto dto.TagListDto
	if err := ctx.ShouldBindQuery(&tagListDto); err != nil {
		zap.L().Error("Failed to bind list tag query", zap.Error(err))
		resp.ParamFail(ctx, err.Error())
		return
	}
	if err := tagListDto.ValidateAndDefault(); err != nil {
		zap.L().Error("Failed to check and format tag list request", zap.Error(err))
		resp.ParamFail(ctx, err.Error())
		return
	}
	tagList, listErr := service.TagService.ListTag(ctx, tagListDto)
	if listErr != nil {
		resp.Fail(ctx, listErr)
		return
	}
	resp.OK(ctx, tagList)
}

func (c *tagHandler) GetTagDetail(ctx *gin.Context) {
	var tagQueryDto dto.TagQueryDto
	if err := ctx.ShouldBindQuery(&tagQueryDto); err != nil {
		zap.L().Error("Failed to bind get tag detail query", zap.Error(err))
		resp.ParamFail(ctx, err.Error())
		return
	}
	if err := tagQueryDto.ValidateAndDefault(); err != nil {
		zap.L().Error("Failed to check and format tag detail request", zap.Error(err))
		resp.ParamFail(ctx, err.Error())
		return
	}
	tagDetail, getErr := service.TagService.GetTagDetail(ctx, tagQueryDto)
	if getErr != nil {
		resp.Fail(ctx, getErr)
		return
	}
	resp.OK(ctx, tagDetail)
}

func (c *tagHandler) CreateTagList(ctx *gin.Context) {
	var tagRequest dto.TagDto
	if err := ctx.ShouldBindJSON(&tagRequest); err != nil {
		zap.L().Error("Failed to bind create tag JSON", zap.Error(err))
		resp.ParamFail(ctx, err.Error())
		return
	}
	if err := tagRequest.ValidateAndDefault(); err != nil {
		zap.L().Error("Failed to check and format tag request", zap.Error(err))
		resp.ParamFail(ctx, err.Error())
		return
	}

	if err := service.TagService.CreateTagList(ctx, tagRequest); err != nil {
		resp.Fail(ctx, err)
		return
	}
	resp.OK(ctx, nil)
}

func (c *tagHandler) UpdateTag(ctx *gin.Context) {
	var tagUpdateDto dto.TagUpdateDto
	if err := ctx.ShouldBindJSON(&tagUpdateDto); err != nil {
		zap.L().Error("Failed to bind update tag JSON", zap.Error(err))
		resp.ParamFail(ctx, err.Error())
		return
	}
	if err := tagUpdateDto.ValidateAndDefault(); err != nil {
		zap.L().Error("Failed to check and format tag request", zap.Error(err))
		resp.ParamFail(ctx, err.Error())
		return
	}
	if err := service.TagService.UpdateTag(ctx, tagUpdateDto); err != nil {
		resp.Fail(ctx, err)
		return
	}
	resp.OK(ctx, nil)
}

func (c *tagHandler) DeleteTagList(ctx *gin.Context) {
	var deleteRequest dto.TagDto
	if err := ctx.ShouldBindJSON(&deleteRequest); err != nil {
		zap.L().Error("Failed to bind delete tag JSON", zap.Error(err))
		resp.ParamFail(ctx, err.Error())
		return
	}
	if err := service.TagService.DeleteTagList(deleteRequest); err != nil {
		resp.Fail(ctx, err)
		return
	}
	resp.OK(ctx, nil)
}
