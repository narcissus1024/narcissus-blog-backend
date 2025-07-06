package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/narcissus1949/narcissus-blog/pkg/dto"
	"github.com/narcissus1949/narcissus-blog/pkg/server/service"
	resp "github.com/narcissus1949/narcissus-blog/pkg/vo/response"
	"go.uber.org/zap"
)

var CategoryController = new(categoryController)

type categoryController struct {
}

func (c *categoryController) ListAllCategory(ctx *gin.Context) {
	categoryList, listErr := service.CategoryService.ListAllCategory(ctx)
	if listErr != nil {
		resp.Fail(ctx, listErr)
		return
	}

	resp.OK(ctx, categoryList)
}

func (c *categoryController) ListCategory(ctx *gin.Context) {
	var categoryDto dto.CategoryListDto
	if err := ctx.ShouldBindQuery(&categoryDto); err != nil {
		zap.L().Error("Failed to bind category list JSON", zap.Error(err))
		resp.ParamFail(ctx, err.Error())
		return
	}
	if err := categoryDto.ValidateAndDefault(); err != nil {
		zap.L().Error("Failed to validate category list request", zap.Error(err))
		resp.ParamFail(ctx, err.Error())
		return
	}
	categoryList, listErr := service.CategoryService.ListCategory(ctx, categoryDto)
	if listErr != nil {
		resp.Fail(ctx, listErr)
		return
	}

	resp.OK(ctx, categoryList)
}

func (c *categoryController) GetCategoryDetail(ctx *gin.Context) {
	var categoryDto dto.CategoryQueryDto
	if err := ctx.ShouldBindQuery(&categoryDto); err != nil {
		zap.L().Error("Failed to bind get category detail JSON", zap.Error(err))
		resp.ParamFail(ctx, err.Error())
		return
	}
	if err := categoryDto.ValidateAndDefault(); err != nil {
		zap.L().Error("Failed to validate get category detail request", zap.Error(err))
		resp.ParamFail(ctx, err.Error())
		return
	}
	categoryDetail, getErr := service.CategoryService.GetCategoryDetail(ctx, categoryDto)
	if getErr != nil {
		resp.Fail(ctx, getErr)
		return
	}
	resp.OK(ctx, categoryDetail)
}

func (c *categoryController) CreateCategoryList(ctx *gin.Context) {
	var categoryDto dto.CategoryDto
	if err := ctx.ShouldBindJSON(&categoryDto); err != nil {
		zap.L().Error("Failed to bind create category JSON", zap.Error(err))
		resp.ParamFail(ctx, err.Error())
		return
	}
	if err := categoryDto.ValidateAndDefault(); err != nil {
		zap.L().Error("Failed to validate create category request", zap.Error(err))
		resp.ParamFail(ctx, err.Error())
		return
	}

	if err := service.CategoryService.CreateCategoryList(ctx, categoryDto); err != nil {

		resp.Fail(ctx, err)
		return
	}
	resp.OK(ctx, nil)
}

func (c *categoryController) UpdateCategory(ctx *gin.Context) {
	var updateRequest dto.CategoryUpdateDto
	if err := ctx.ShouldBindJSON(&updateRequest); err != nil {
		zap.L().Error("Failed to bind update category JSON", zap.Error(err))
		resp.ParamFail(ctx, err.Error())
		return
	}
	if err := updateRequest.ValidateAndDefault(); err != nil {
		zap.L().Error("Failed to validate update category request", zap.Error(err))
		resp.ParamFail(ctx, err.Error())
		return
	}

	if err := service.CategoryService.UpdateCategory(ctx, updateRequest); err != nil {
		resp.Fail(ctx, err)
		return
	}
	resp.OK(ctx, nil)
}

func (c *categoryController) DeleteCategoryList(ctx *gin.Context) {
	var deleteRequest dto.CategoryDto
	if err := ctx.ShouldBindJSON(&deleteRequest); err != nil {
		zap.L().Error("Failed to bind delete category JSON", zap.Error(err))
		resp.ParamFail(ctx, err.Error())
		return
	}
	if err := deleteRequest.ValidateAndDefault(); err != nil {
		zap.L().Error("Failed to validate delete category request", zap.Error(err))
		resp.ParamFail(ctx, err.Error())
		return
	}

	if err := service.CategoryService.DeleteCategoryList(ctx, deleteRequest); err != nil {

		resp.Fail(ctx, err)
		return
	}
	resp.OK(ctx, nil)
}
