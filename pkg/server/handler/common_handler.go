package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/narcissus1949/narcissus-blog/internal/logger"
	"github.com/narcissus1949/narcissus-blog/pkg/dto"
	"github.com/narcissus1949/narcissus-blog/pkg/server/service"
	resp "github.com/narcissus1949/narcissus-blog/pkg/vo/response"
	"go.uber.org/zap"
)

var CommonHandler = new(commonHandler)

type commonHandler struct {
}

func (c *commonHandler) UploadImage(ctx *gin.Context) {
	// 限制本次请求体最大为 2MB
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, 2<<20)

	file, err := ctx.FormFile("file")
	if err != nil {
		logger.FromContext(ctx.Request.Context()).Error("Failed to get image file from request", zap.Error(err))
		resp.ParamFail(ctx, err.Error())
		return
	}

	result, uploadErr := service.CommonServiceInstance.UploadImage(ctx, file)
	if uploadErr != nil {
		resp.Fail(ctx, uploadErr)
		return
	}

	resp.OK(ctx, result)
}

func (c *commonHandler) GetRASPublicKey(ctx *gin.Context) {
	result, err := service.CommonServiceInstance.GetRASPublicKey(ctx)
	if err != nil {
		resp.Fail(ctx, err)
		return
	}
	resp.OK(ctx, result)
}

func (c *commonHandler) PublicKeyEncrypt(ctx *gin.Context) {
	var req dto.PublicKeyEncrypDto
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.FromContext(ctx.Request.Context()).Error("Failed to bind public key encrypt request", zap.Error(err))
		resp.ParamFail(ctx, err.Error())
		return
	}
	result, err := service.CommonServiceInstance.PublicKeyEncrypt(req)
	if err != nil {
		resp.Fail(ctx, err)
		return
	}
	resp.OK(ctx, result)
}
