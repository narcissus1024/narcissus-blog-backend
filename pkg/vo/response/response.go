package resp

import (
	"net/http"

	"github.com/gin-gonic/gin"
	cerr "github.com/narcissus1949/narcissus-blog/internal/error"
	"github.com/narcissus1949/narcissus-blog/pkg/vo/result"
)

func ResponseJson(c *gin.Context, httpCode int, result result.Result) {
	c.JSON(httpCode, result)
}

func Fail(c *gin.Context, err error) {
	ResponseJson(c, http.StatusOK, result.Fail(c, err))
}

func ParamFail(c *gin.Context, msg ...string) {
	ResponseJson(c, http.StatusBadRequest, result.Fail(c, cerr.NewParamError(msg...)))
}

func UnauthorizedFail(c *gin.Context, msg ...string) {
	ResponseJson(c, http.StatusUnauthorized, result.Fail(c, cerr.New(cerr.UNAUTHORIZED, msg...)))
}

func TokenExpire(c *gin.Context) {
	ResponseJson(c, http.StatusUnauthorized, result.Fail(c, cerr.New(cerr.ERROR_USER_TOKEN_EXPIRE)))
}
func OK(c *gin.Context, data any) {
	ResponseJson(c, http.StatusOK, result.Success(c, data))
}
