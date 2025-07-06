package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/narcissus1949/narcissus-blog/internal/database/cache"
	"github.com/narcissus1949/narcissus-blog/internal/jwt"
	"github.com/narcissus1949/narcissus-blog/internal/utils"
	resp "github.com/narcissus1949/narcissus-blog/pkg/vo/response"
	"go.uber.org/zap"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 存放token一般都是在请求头部的Authorization按该格式存放"Bearer [token]"
		authorization := c.Request.Header.Get("Authorization")
		if authorization == "" {
			zap.L().Error("Authorization is empty")
			resp.UnauthorizedFail(c)
			c.Abort()
			return
		}

		token, parseAuthorizationErr := utils.ParseAuthorization(authorization)
		if parseAuthorizationErr != nil {
			zap.L().Error("Failed to parse authorization", zap.Error(parseAuthorizationErr))
			resp.UnauthorizedFail(c)
			c.Abort()
			return
		}

		// 解析token
		claims, parseTokenErr := jwt.ParseToken(token)
		if parseTokenErr != nil {
			if jwt.IsExpireErr(parseTokenErr) {
				zap.L().Error("Token is expired", zap.Error(parseTokenErr))
				resp.TokenExpire(c)
				c.Abort()
				return
			}
			resp.UnauthorizedFail(c)
			c.Abort()
			return
		}

		// 验证token是否在黑名单
		redisClient := cache.Client
		count, tokenExistsErr := redisClient.Exists(c, utils.ACCESS_TOKEN_BLACKLIST+token, utils.REFRESH_TOKEN_BLACKLIST+token).Result()
		if tokenExistsErr != nil {
			zap.L().Error("Failed to check token exists", zap.Error(tokenExistsErr))
			resp.UnauthorizedFail(c)
			c.Abort()
			return
		}
		if count > 0 {
			// token无效，可能由登出/注销导致
			zap.L().Error("Token is invalid, in blacklist")
			resp.UnauthorizedFail(c)
			c.Abort()
			return
		}

		// 向context中存放信息
		c.Set(utils.CONTEXT_USER_ID, claims.UserID)

		c.Next()
	}
}
