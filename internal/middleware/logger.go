package middleware

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/narcissus1949/narcissus-blog/internal/logger"

	"github.com/gin-gonic/gin"
	"github.com/narcissus1949/narcissus-blog/internal/utils"
	"go.uber.org/zap"
)

// GinLogger 接收gin框架默认的日志
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		// x-request-id
		requestID := c.GetHeader(utils.X_REQUEST_ID)
		if requestID == "" {
			requestID = utils.GenerateUUID()
		}
		c.Set(utils.X_REQUEST_ID, requestID)
		c.Writer.Header().Set(utils.X_REQUEST_ID, requestID)

		// 为该请求构建带有 X-Request-ID 的 logger，并注入到 request context
		reqLogger := zap.L().With(zap.String(utils.X_REQUEST_ID, requestID))
		c.Request = c.Request.WithContext(logger.ToContext(c.Request.Context(), reqLogger))
		// body, _ := c.GetRawData()
		// c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		c.Next()

		cost := time.Since(start)
		reqLogger.Info("http request info",
			zap.String(utils.X_REQUEST_ID, requestID),
			zap.String("path", path),
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			// zap.String("body", string(body)),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

// GinRecovery recover掉项目可能出现的panic，并使用zap记录相关日志
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 从 ctx 中获取请求级 logger（若没有则回退到全局）
				reqLogger := logger.FromContext(c.Request.Context())
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					reqLogger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					reqLogger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					reqLogger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
