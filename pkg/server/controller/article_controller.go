package controller

import (
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/narcissus1949/narcissus-blog/pkg/dto"
	"github.com/narcissus1949/narcissus-blog/pkg/server/service"
	resp "github.com/narcissus1949/narcissus-blog/pkg/vo/response"
	"go.uber.org/zap"
)

var ArticleController = new(articleController)

type articleController struct {
}

func (c *articleController) SaveArticle(ctx *gin.Context) {
	var articleDto dto.ArticleDto
	if err := ctx.ShouldBindJSON(&articleDto); err != nil {
		zap.L().Error("Failed to bind article JSON", zap.Error(err))
		resp.ParamFail(ctx, err.Error())
		return
	}
	if err := articleDto.VlidateAndDefault(); err != nil {
		zap.L().Error("Failed to validate article request", zap.Error(err))
		resp.ParamFail(ctx, err.Error())
		return
	}

	if articleDto.ID == nil || *articleDto.ID <= 0 {
		if err := service.ArticleService.CreateArticle(ctx, articleDto); err != nil {
			resp.Fail(ctx, err)
			return
		}
	} else {
		if err := service.ArticleService.UpdateArticle(ctx, articleDto); err != nil {
			resp.Fail(ctx, err)
			return
		}
	}

	resp.OK(ctx, nil)
}

func (c *articleController) CreateArticle(ctx *gin.Context) {
	var articleDto dto.ArticleDto
	if err := ctx.ShouldBindJSON(&articleDto); err != nil {
		zap.L().Error("Failed to bind article JSON", zap.Error(err))
		resp.ParamFail(ctx, err.Error())
		return
	}
	if err := articleDto.VlidateAndDefault(); err != nil {
		zap.L().Error("Failed to validate article request", zap.Error(err))
		resp.ParamFail(ctx, err.Error())
		return
	}
	if err := service.ArticleService.CreateArticle(ctx, articleDto); err != nil {
		resp.Fail(ctx, err)
		return
	}

	resp.OK(ctx, nil)
}

func (c *articleController) ListArticleAdmin(ctx *gin.Context) {
	var articleListRequest dto.ArticleListDto
	// 先设置默认，防止binding校验失败
	if err := articleListRequest.VlidateAndSetDefault(); err != nil {
		zap.L().Error("Failed to validate article list request", zap.Error(err))
		resp.ParamFail(ctx, err.Error())
		return
	}
	if err := ctx.ShouldBindJSON(&articleListRequest); err != nil && !errors.Is(err, io.EOF) {
		zap.L().Error("Failed to bind article list JSON", zap.Error(err))
		resp.ParamFail(ctx, err.Error())
		return
	}
	// 第二次设置默认，防止传入零值
	if err := articleListRequest.VlidateAndSetDefault(); err != nil {
		zap.L().Error("Failed to validate article list request", zap.Error(err))
		resp.ParamFail(ctx, err.Error())
		return
	}

	articleList, err := service.ArticleService.ListArticleAdmin(ctx, articleListRequest)
	if err != nil {
		resp.Fail(ctx, err)
		return
	}
	resp.OK(ctx, articleList)
}

func (c *articleController) ListArticle(ctx *gin.Context) {
	var articleListRequest dto.ArticleListDto
	// 先设置默认，防止binding校验失败
	if err := articleListRequest.VlidateAndSetDefault(); err != nil {
		zap.L().Error("Failed to validate article list request", zap.Error(err))
		resp.ParamFail(ctx, err.Error())
		return
	}
	if err := ctx.ShouldBindJSON(&articleListRequest); err != nil && !errors.Is(err, io.EOF) {
		zap.L().Error("Failed to bind article list JSON", zap.Error(err))
		resp.ParamFail(ctx, err.Error())
		return
	}
	// 第二次设置默认，防止传入零值
	if err := articleListRequest.VlidateAndSetDefault(); err != nil {
		zap.L().Error("Failed to validate article list request", zap.Error(err))
		resp.ParamFail(ctx, err.Error())
		return
	}

	// 只允许查已上线的文章
	status := true
	articleListRequest.Status = &status
	articleList, err := service.ArticleService.ListArticleAdmin(ctx, articleListRequest)
	if err != nil {
		resp.Fail(ctx, err)
		return
	}
	resp.OK(ctx, articleList)
}

func (c *articleController) GetArticleeDetail(ctx *gin.Context) {
	idStr := ctx.Query("id")
	if len(strings.TrimSpace(idStr)) == 0 {
		zap.L().Error("Failed to get article detail, id is 0")
		resp.ParamFail(ctx, errors.New("id invalide").Error())
		return
	}
	id, idConvErr := strconv.Atoi(idStr)
	if idConvErr != nil {
		zap.L().Error("Failed to convert article id to int", zap.Error(idConvErr))
		resp.ParamFail(ctx, idConvErr.Error())
		return
	}

	articleDetail, getDetailErr := service.ArticleService.GetArticleDetail(ctx, int64(id))
	if getDetailErr != nil {
		resp.Fail(ctx, getDetailErr)
		return
	}
	resp.OK(ctx, articleDetail)
}

func (c *articleController) DeleteArticleList(ctx *gin.Context) {
	var deleteDto dto.ArticleDeleteDto
	if err := ctx.ShouldBindJSON(&deleteDto); err != nil {
		zap.L().Error("Failed to bind delete article list JSON", zap.Error(err))
		resp.ParamFail(ctx, err.Error())
		return
	}
	if err := deleteDto.VlidateAndDefault(); err != nil {
		zap.L().Error("Failed to validate delete article list request", zap.Error(err))
		resp.ParamFail(ctx, err.Error())
		return
	}
	if err := service.ArticleService.DeleteArticleList(ctx, deleteDto); err != nil {
		resp.Fail(ctx, err)
		return
	}
	resp.OK(ctx, nil)
}
