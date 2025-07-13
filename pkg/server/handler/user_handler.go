package handler

import (
	"github.com/gin-gonic/gin"
	cerr "github.com/narcissus1949/narcissus-blog/internal/error"
	"github.com/narcissus1949/narcissus-blog/pkg/dto"
	"github.com/narcissus1949/narcissus-blog/pkg/server/service"
	resp "github.com/narcissus1949/narcissus-blog/pkg/vo/response"
	"go.uber.org/zap"
)

var UserHandler = new(userHandler)

type userHandler struct {
	// service service.UserService
}

func (c *userHandler) SignIn(ctx *gin.Context) {
	var signin dto.SignIn
	if err := ctx.ShouldBindJSON(&signin); err != nil {
		zap.L().Error("Failed to bind signin JSON", zap.Error(err))
		resp.ParamFail(ctx, err.Error())
		return
	}

	if err := service.UserServiceInstance.SignIn(signin); err != nil {
		resp.Fail(ctx, err)
		return
	}

	resp.OK(ctx, nil)
}

func (c *userHandler) Login(ctx *gin.Context) {
	var loginDto dto.LoginDto
	if err := ctx.ShouldBindJSON(&loginDto); err != nil {
		zap.L().Error("Failed to bind login JSON", zap.Error(err))
		resp.ParamFail(ctx, err.Error())
		return
	}
	token, err := service.UserServiceInstance.Login(loginDto)
	if err != nil {
		resp.Fail(ctx, err)
		return
	}
	if token == nil {
		resp.Fail(ctx, cerr.NewSysError("登陆失败"))
		return
	}

	resp.OK(ctx, token)
}

func (c *userHandler) Logout(ctx *gin.Context) {
	var logoutDto dto.LogoutDto
	if err := ctx.ShouldBindJSON(&logoutDto); err != nil {
		zap.L().Error("Failed to bind logout JSON", zap.Error(err))
		resp.ParamFail(ctx, err.Error())
		return
	}
	if err := logoutDto.Validate(); err != nil {
		zap.L().Error("Failed to validate logout request", zap.Error(err))
		resp.ParamFail(ctx, err.Error())
		return
	}

	err := service.UserServiceInstance.Logout(ctx, logoutDto)
	if err != nil {
		resp.Fail(ctx, err)
		return
	}

	resp.OK(ctx, nil)
}

func (c *userHandler) RefreshToken(ctx *gin.Context) {
	var refreshTokenDto dto.RefreshTokenDto
	if err := ctx.ShouldBindJSON(&refreshTokenDto); err != nil {
		zap.L().Error("Failed to bind refresh token JSON", zap.Error(err))
		resp.ParamFail(ctx, err.Error())
		return
	}
	if err := refreshTokenDto.Validate(); err != nil {
		zap.L().Error("Failed to validate refresh token request", zap.Error(err))
		resp.ParamFail(ctx, err.Error())
		return
	}

	result, refreshTokenErr := service.UserServiceInstance.RefreshToken(ctx, refreshTokenDto)
	if refreshTokenErr != nil {
		resp.Fail(ctx, refreshTokenErr)
		return
	}
	resp.OK(ctx, result)
}
