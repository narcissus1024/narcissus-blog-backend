package result

import (
	"errors"

	"github.com/gin-gonic/gin"
	cerr "github.com/narcissus1949/narcissus-blog/internal/error"
)

type Result struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      any         `json:"data"`
	RequestId interface{} `json:"requestId"`
}

func Success(c *gin.Context, data any) Result {
	return Result{
		Code:    cerr.SUCCESS,
		Message: cerr.GetMessage(cerr.SUCCESS),
		Data:    data,
	}
}

func Fail(c *gin.Context, err error) Result {
	var cerrErr *cerr.Error
	if ok := errors.As(err, &cerrErr); ok {
		return Result{
			Code:    cerrErr.Code,
			Message: cerrErr.Message,
		}
	}

	return Result{
		Code:    cerr.SYSTEM_ERROR,
		Message: cerr.GetMessage(cerr.SYSTEM_ERROR),
	}
}
