package handler

import (
	"github.com/gin-gonic/gin"
	resp "github.com/narcissus1949/narcissus-blog/pkg/vo/response"
)

type healthHandler struct{}

var HealthHandler = new(healthHandler)

// Ping 健康检查接口
func (h *healthHandler) Ping(c *gin.Context) {
	resp.OK(c, "pong")
}
