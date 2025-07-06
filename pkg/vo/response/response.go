package resp

import (
	"net/http"

	"github.com/gin-gonic/gin"
	cerr "github.com/narcissus1949/narcissus-blog/internal/error"
	"github.com/narcissus1949/narcissus-blog/pkg/vo"
)

func ResponseJson(c *gin.Context, httpCode int, err *cerr.Error, data any) {
	result := vo.Result{
		Error: err,
		Data:  data,
	}
	c.JSON(httpCode, result)
}

func OK(c *gin.Context, data any) {
	ResponseJson(c, http.StatusOK, cerr.NewSuccess(), data)
}

func Fail(c *gin.Context, err *cerr.Error) {
	ResponseJson(c, http.StatusOK, err, nil)
}

func ParamFail(c *gin.Context, msg ...string) {
	ResponseJson(c, http.StatusBadRequest, cerr.NewParamError(msg...), nil)
}

func UnauthorizedFail(c *gin.Context, msg ...string) {
	ResponseJson(c, http.StatusUnauthorized, cerr.New(cerr.UNAUTHORIZED, msg...), nil)
}

func TokenExpire(c *gin.Context) {
	ResponseJson(c, http.StatusUnauthorized, cerr.New(cerr.ERROR_USER_TOKEN_EXPIRE), nil)
}
